package domain

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	"upper.io/db.v3"

	"github.com/star-table/startable-server/app/service/projectsvc/domain/lc_pro_domain"

	tableV1 "gitea.bjx.cloud/LessCode/interface/golang/table/v1"
	"github.com/star-table/startable-server/app/facade/appfacade"
	"github.com/star-table/startable-server/app/facade/idfacade"
	"github.com/star-table/startable-server/app/facade/orgfacade"
	"github.com/star-table/startable-server/app/facade/permissionfacade"
	"github.com/star-table/startable-server/app/facade/resourcefacade"
	"github.com/star-table/startable-server/app/facade/tablefacade"
	"github.com/star-table/startable-server/app/service/projectsvc/dao"
	"github.com/star-table/startable-server/app/service/projectsvc/po"
	"github.com/star-table/startable-server/common/core/consts"
	"github.com/star-table/startable-server/common/core/errs"
	"github.com/star-table/startable-server/common/core/util"
	"github.com/star-table/startable-server/common/core/util/copyer"
	"github.com/star-table/startable-server/common/core/util/date"
	"github.com/star-table/startable-server/common/core/util/format"
	"github.com/star-table/startable-server/common/core/util/json"
	"github.com/star-table/startable-server/common/core/util/slice"
	"github.com/star-table/startable-server/common/extra/lc_helper"
	"github.com/star-table/startable-server/common/library/db/mysql"
	"github.com/star-table/startable-server/common/model/bo"
	"github.com/star-table/startable-server/common/model/vo"
	"github.com/star-table/startable-server/common/model/vo/appvo"
	"github.com/star-table/startable-server/common/model/vo/lc_table"
	"github.com/star-table/startable-server/common/model/vo/orgvo"
	"github.com/star-table/startable-server/common/model/vo/permissionvo"
	"github.com/star-table/startable-server/common/model/vo/projectvo"
	"github.com/star-table/startable-server/common/model/vo/resourcevo"
	"google.golang.org/protobuf/types/known/structpb"
	"upper.io/db.v3/lib/sqlbuilder"
)

func CreateProject(entity bo.ProjectBo, orgId, currentUserId int64, input vo.CreateProjectReq, remarkObj orgvo.OrgRemarkConfigType) (int64, bo.ProjectBo, []int64, errs.SystemErrorInfo) {
	appId := int64(0)

	memberEntities, addedMemberIds, err := HandleProjectMember(orgId, currentUserId, entity.Owner, entity.Id, input.MemberIds, input.FollowerIds, input.IsAllMember, input.MemberForDepartmentID, input.OwnerIds)
	if err != nil {
		log.Errorf("[CreateProject] HandleProjectMember err: %v", err)
		return appId, entity, nil, err
	}

	//查询资源是否已存在
	var resourceId int64

	dealResourcePath(&entity, orgId, input.ResourcePath, input.ResourceType, &resourceId)

	err1 := mysql.TransX(func(tx sqlbuilder.Tx) error {
		//插入项目成员
		if insertErr := insertMemberEntities(tx, memberEntities); insertErr != nil {
			log.Errorf("[CreateProject] insertMemberEntities err: %v", insertErr)
			return insertErr
		}

		//插入资源
		err = insertSource(&entity, input.ResourcePath, resourceId, tx, orgId, currentUserId, input.ResourceType)
		if err != nil {
			log.Errorf("[CreateProject] insertSource err: %v", err)
			return errs.BuildSystemErrorInfo(errs.MysqlOperateError, err)
		}

		//创建无码应用
		configJson := getProjectColumnConfig(*input.ProjectTypeID)
		appResp, err := CreateAppInLessCode(orgId, currentUserId, input.Name, remarkObj.OrgSummaryTableAppId,
			input.ParentID, entity.Id, configJson, int(*input.ProjectTypeID))
		if err != nil {
			log.Errorf("[CreateProject] CreateAppInLessCode err: %v", err)
			return err
		}
		appId = appResp.Id
		entity.AppId = appId

		// 无码交互：创建应用后，更新定制化的应用权限组的权限项
		if appId > 0 {
			if err := lc_pro_domain.UpdateOpForAppPermissionGroup(appId); err != nil {
				log.Errorf("[CreateProject] lc_pro_domain.UpdateOpForAppPermissionGroup err: %v", err)
				return err
			}
		} else {
			// 调用无码系统，创建应用如果不成功，则不能继续执行。
			err := errs.LcUpdateAppPermissionGroupOptAuth
			log.Errorf("[CreateProject] err: %v", err)
			return err
		}

		//插入项目
		projectError := insertProject(tx, entity)
		if projectError != nil {
			log.Errorf("[CreateProject] insertProject err: %v", projectError)
			return projectError
		}

		//插入项目公告、日历配置、群聊开启配置等数据
		syncCalendarFlag := TransferSyncOutCalendarStatusIntoOne(input.SyncCalendarStatusList)
		// 项目隐私模式。默认不开启2。隐私模式已下线。这里只是一个默认值，不用管。
		defaultProPrivacy := consts.ProSetPrivacyDisable
		if input.PrivacyStatus == nil {
			input.PrivacyStatus = &defaultProPrivacy
		}

		err = insertProjectDetail(&entity, orgId, currentUserId, &syncCalendarFlag, input.IsCreateFsChat, input.PrivacyStatus)
		if err != nil {
			log.Errorf("[CreateProject] insertProjectDetail err: %v", err)
			return err
		}

		return nil
	})
	if err1 != nil {
		log.Errorf("tx.Commit() err: %v", err1)
		return appId, entity, nil, errs.BuildSystemErrorInfo(errs.MysqlOperateError, err1)
	}

	return appId, entity, addedMemberIds, nil
}

func getProjectColumnConfig(projectTypeId int64) string {
	notNeedSummeryColumnIds := GetNoNeedColumnByProjectType(projectTypeId)
	// 创建不同类型的项目，对应的表不允许有一些特定的列
	config := map[string]interface{}{
		"notNeedSummeryColumnIds": notNeedSummeryColumnIds,
	}

	// 如果是空应用，不需要优先级
	if projectTypeId != consts.ProjectTypeEmpty {
		config["fields"] = []interface{}{lc_helper.GetDocumentColumn(), lc_helper.GetOrgPriorityField()}
	} else {
		config["fields"] = []interface{}{lc_helper.GetSelectColumn(), lc_helper.GetMultiSelectColumn(), lc_helper.GetDocumentColumn()}
	}

	return json.ToJsonIgnoreError(config)
}

// CreateAppInLessCode 创建项目对应的应用。
// notNeedSummeryColumnIds 创建不同类型的项目，对应的表不允许有一些特定的列
func CreateAppInLessCode(orgId, opUserId int64, proName string, extendsId int64, lcFolderId *int64, projectId int64,
	config string, projectType int,
) (*permissionvo.CreateLessCodeAppRespData, errs.SystemErrorInfo) {
	// 4表示极星项目
	appType := 4
	parentId := int64(0)
	if lcFolderId != nil {
		parentId = *lcFolderId
	}
	req := permissionvo.CreateLessCodeAppReq{
		OrgId:       &orgId,
		AppType:     &appType,
		Name:        &proName,
		UserId:      &opUserId,
		Config:      config, // GetAppDefaultColumnConfig() 列信息和汇总表基本一样，可以不传。如果有特殊定制，可以传入配置 json
		ProjectId:   projectId,
		ExtendsId:   extendsId,
		ParentId:    parentId,
		ProjectType: projectType,
	}
	resp := appfacade.CreateLessCodeApp(&req)
	if resp.Failure() {
		log.Error(resp.Error())
		return nil, resp.Error()
	}

	return resp.Data, nil
}

// GetDefaultColumnForTaskBar 通用项目，默认的任务栏，taskBar <==> projectObjectType。select 类型
func GetDefaultColumnForTaskBar() lc_table.LcCtSelect {
	column := lc_helper.GetLcCtSelect(consts.BasicFieldProjectObjectTypeId, "任务栏", "Task Bar",
		"select", lc_helper.GetDefaultSelectOptionsForTaskBar(), false, true, true, false, false)
	return column
}

// GetDefaultColumnForIssueStatus 通用项目，默认的任务状态 option list。groupSelect 类型
//func GetDefaultColumnForIssueStatus() lc_table.LcOneColumn {
//	column := lc_helper.GetLcCtGroupSelect(consts.BasicFieldIssueStatus, "任务状态", "groupSelect",
//		lc_helper.GetDefaultGroupSelectForIssueStatus(), true)
//	return column
//}

func dealResourcePath(entity *bo.ProjectBo, orgId int64, resourcePath string, resourceType int, resourceId *int64) {
	if resourcePath != "" {
		respVo := resourcefacade.GetIdByPath(
			resourcevo.GetIdByPathReqVo{
				OrgId:        orgId,
				ResourceType: resourceType,
				ResourcePath: resourcePath,
			})
		if !respVo.Failure() {
			*resourceId = respVo.ResourceId
			(*entity).ResourceId = *resourceId
		} else {
			//log.Errorf("[dealResourcePath] orgId:%v,resourceType:%v,resourcePath:%v, err: %v", orgId, resourceType, resourcePath, respVo.Error())
		}
	}
}

func insertProjectDetail(entity *bo.ProjectBo, orgId, currentUserId int64, syncCalendarFlag, isCreateFsChat, proPrivacyStatus *int) errs.SystemErrorInfo {

	detailId, err := idfacade.ApplyPrimaryIdRelaxed(consts.TableProjectDetail)
	if err != nil {
		return errs.BuildSystemErrorInfo(errs.ApplyIdError, err)
	}
	// IsSyncOutCalendar := consts.IsSyncOutCalendar
	IsNotSyncOutCalendar := consts.IsNotSyncOutCalendar
	if syncCalendarFlag != nil && *syncCalendarFlag > 0 {
		// 无需赋值，只需用传入的参数值。
		// syncCalendarFlag = &IsSyncOutCalendar
	} else {
		syncCalendarFlag = &IsNotSyncOutCalendar
	}
	insertProjectDetailErr := dao.InsertProjectDetail(po.PpmProProjectDetail{
		Id:                detailId,
		OrgId:             orgId,
		ProjectId:         entity.Id,
		Notice:            consts.BlankString,
		IsSyncOutCalendar: *syncCalendarFlag,
		Creator:           currentUserId,
		CreateTime:        time.Now(),
		Updator:           currentUserId,
		UpdateTime:        time.Now(),
		IsDelete:          consts.AppIsNoDelete,
	})
	if insertProjectDetailErr != nil {
		return errs.BuildSystemErrorInfo(errs.MysqlOperateError, insertProjectDetailErr)
	}

	return nil
}

func insertProject(tx sqlbuilder.Tx, entity bo.ProjectBo) errs.SystemErrorInfo {

	insert := &po.PpmProProject{}
	err1 := copyer.Copy(entity, insert)
	if err1 != nil {
		return errs.BuildSystemErrorInfo(errs.SystemError, err1)
	}
	_, insertProjectErr := tx.Collection(consts.TableProject).Insert(insert)
	if insertProjectErr != nil {
		return errs.BuildSystemErrorInfo(errs.MysqlOperateError, insertProjectErr)
	}
	return nil
}

func insertMemberEntities(tx sqlbuilder.Tx, memberEntities []interface{}) error {
	if len(memberEntities) == 0 {
		return nil
	}

	err := PaginationInsert(memberEntities, &po.PpmProProjectRelation{}, tx)
	if err != nil {
		log.Errorf("[insertMemberEntities] err: %v", err)
		return errs.BuildSystemErrorInfo(errs.MysqlOperateError, err)
	}

	return nil
}

func insertSource(entity *bo.ProjectBo, resourcePath string, resourceId int64, tx sqlbuilder.Tx, orgId, currentUserId int64,
	resourceType int) errs.SystemErrorInfo {
	if resourcePath != "" && resourceId == 0 {
		fileName := util.ParseFileName(resourcePath)
		suffix := util.ParseFileSuffix(fileName)
		respVo := resourcefacade.CreateResource(resourcevo.CreateResourceReqVo{
			CreateResourceBo: bo.CreateResourceBo{
				Path:       resourcePath,
				Name:       fileName,
				Suffix:     suffix,
				OrgId:      orgId,
				OperatorId: currentUserId,
				Type:       resourceType,
			},
		})
		if respVo.Failure() {
			return respVo.Error()
		}
		entity.ResourceId = respVo.ResourceId
	}
	return nil
}

// 判断项目名是否重复
func JudgeRepeatProjectName(name *string, orgId int64, projectId *int64) (string, errs.SystemErrorInfo) {
	if name == nil {
		*name = consts.BlankString
	}
	cond := make(db.Cond)
	cond = db.Cond{
		consts.TcIsDelete: db.Eq(consts.AppIsNoDelete),
		consts.TcName:     db.Eq(name),
		consts.TcOrgId:    orgId,
	}
	//如果传项目id
	if projectId != nil {
		cond[consts.TcId] = db.NotEq(projectId)
	}
	exist, err := mysql.IsExistByCond(consts.TableProject, cond)
	if err != nil {
		return *name, errs.BuildSystemErrorInfo(errs.MysqlOperateError, err)
	}
	if exist {
		return *name, errs.BuildSystemErrorInfo(errs.RepeatProjectName)
	}

	return *name, nil
}

// 判断前缀编号是否重复
func JudgeRepeatProjectPreCode(preCode *string, orgId int64, projectId *int64) (string, errs.SystemErrorInfo) {
	if preCode == nil {
		*preCode = consts.BlankString
	}
	cond := make(db.Cond)
	cond = db.Cond{
		consts.TcIsDelete: db.Eq(consts.AppIsNoDelete),
		consts.TcPreCode:  preCode,
		consts.TcOrgId:    orgId,
	}
	//如果传项目id
	if projectId != nil {
		cond[consts.TcId] = db.NotEq(projectId)
	}
	exist, err := mysql.IsExistByCond(consts.TableProject, cond)
	if err != nil {
		return *preCode, errs.BuildSystemErrorInfo(errs.MysqlOperateError, err)
	}
	if exist {
		return *preCode, errs.BuildSystemErrorInfo(errs.RepeatProjectPrecode)
	}

	return *preCode, nil
}

// GetTypeAndStatus 获取项目类型和初始状态
func GetTypeAndStatus(projectTypeId int64) (int64, int64, errs.SystemErrorInfo) {
	if projectTypeId == 0 {
		projectTypeId = consts.ProjectTypeCommon2022V47
	}
	return projectTypeId, consts.StatusRunning.ID, nil
}

// CreateProjectTables 任务状态改造后，创建项目，会创建对应的 table list
//func CreateProjectTables(orgId, projectId, currentUserId, projectTypeId, proAppId int64, tx sqlbuilder.Tx) errs.SystemErrorInfo {
//	defaultTablesMap := consts.DefaultProTypeMap2TableList
//	defaultTables, ok := defaultTablesMap[projectTypeId]
//	if !ok {
//		err := errs.NoSupportProjectType
//		log.Errorf("[CreateProjectTables] orgId: %d, projectId: %d, err: %v", orgId, projectId, err)
//		return err
//	}
//	// 创建 table list
//	for _, proTable := range defaultTables {
//		if resp := tablefacade.CreateTable(projectvo.CreateTableReq{
//			OrgId:  orgId,
//			UserId: currentUserId,
//			Input: &tableV1.CreateTableRequest{
//				AppId:            proAppId,
//				Name:             proTable.Name,
//				BasicColumns:     nil,
//				IsNeedStoreTable: false,
//				IsNeedColumn:     true,
//				Columns:          nil,
//			},
//		}); resp.Failure() {
//			log.Errorf("[CreateProjectTables] tablefacade.CreateTable err: %v", resp.Error())
//			return resp.Error()
//		}
//	}
//
//	return nil
//}

func validateOrder(order string) (string, bool) {
	re := regexp.MustCompile(" +")
	strArr := re.Split(order, -1)
	if len(strArr) != 2 {
		return "", false
	}
	if ok, _ := slice.Contain([]string{"asc", "desc"}, strings.ToLower(strArr[1])); !ok {
		return "", false
	}
	return fmt.Sprintf("`%s` %s,", strArr[0], strArr[1]), format.VerifySqlFieldFormat(strArr[0])
}

func GetProjectList(currentUserId int64, joinParams db.Cond, union *db.Union, order []*string, size int, page int) ([]*bo.ProjectBo, int64, errs.SystemErrorInfo) {
	conn, err := mysql.GetConnect()
	if err != nil {
		return nil, 0, errs.BuildSystemErrorInfo(errs.MysqlOperateError, err)
	}
	entities := &[]*po.PpmProProject{}
	mid := conn.Collection(consts.TableProject).Find(joinParams)
	if union != nil {
		mid = mid.And(union)
	}
	log.Infof("[GetProjectList] %v", mid.String())
	orderRaw := []string{}

	if order != nil {
		for _, v := range order {
			if v == nil {
				continue
			}

			orderRaw = append(orderRaw, *v)
		}
	}
	if len(orderRaw) > 0 {
		orderRawStr := strings.Join(orderRaw, ",")
		mid = mid.OrderBy(db.Raw(orderRawStr))
	}
	count, err := mid.TotalEntries()
	if err != nil {
		log.Error(err)
		return nil, int64(count), errs.BuildSystemErrorInfo(errs.MysqlOperateError, err)
	}
	if size > 0 && page > 0 {
		err = mid.Paginate(uint(size)).Page(uint(page)).All(entities)
	} else {
		err = mid.All(entities)
	}
	if err != nil {
		return nil, 0, errs.BuildSystemErrorInfo(errs.MysqlOperateError, err)
	}
	projectBos := &[]*bo.ProjectBo{}
	_ = copyer.Copy(entities, projectBos)

	return *projectBos, int64(count), nil
}

// GetProjectSimple 仅仅是获取项目的信息，不包括项目的管理员
func GetProjectSimple(orgId int64, projectId int64, tx ...sqlbuilder.Tx) (*bo.ProjectBo, errs.SystemErrorInfo) {
	project := po.PpmProProject{}
	cond := db.Cond{
		consts.TcIsDelete: consts.AppIsNoDelete,
		consts.TcId:       projectId,
		consts.TcOrgId:    orgId,
	}
	var err error
	if len(tx) > 0 {
		err = mysql.TransSelectOneByCond(tx[0], project.TableName(), cond, &project)
	} else {
		err = mysql.SelectOneByCond(project.TableName(), cond, &project)
	}
	if err != nil {
		log.Errorf("[GetProjectSimple]获取项目表数据异常, orgId:%d, projectId:%d, 错误: %v", orgId, projectId, err)
		return nil, errs.BuildSystemErrorInfo(errs.MysqlOperateError, errs.ProjectNotExist)
	}

	projectBo := bo.ProjectBo{}
	err1 := copyer.Copy(&project, &projectBo)
	if err1 != nil {
		log.Error(err1)
		return nil, errs.BuildSystemErrorInfo(errs.ObjectCopyError)
	}

	return &projectBo, nil
}

// 目前只有在 UpdateProjectWithoutAuth 这里用到了 ownerIds，所以暂时改了个函数名
// 让UpdateProjectWithoutAuth调用，其他 原来调用 GetProjectInfo的地方去掉查询负责人的逻辑
func GetProjectInfoWithOwnerIds(id int64, orgId int64) (bo.ProjectBo, errs.SystemErrorInfo) {
	project := &po.PpmProProject{}
	err := mysql.SelectOneByCond(project.TableName(), db.Cond{
		consts.TcIsDelete: consts.AppIsNoDelete,
		consts.TcId:       id,
		consts.TcOrgId:    orgId,
	}, project)
	projectBo := &bo.ProjectBo{}
	if err != nil {
		return *projectBo, errs.BuildSystemErrorInfo(errs.MysqlOperateError, err)
	}
	_ = copyer.Copy(project, projectBo)
	//负责人集合
	relationPos := &[]po.PpmProProjectRelation{}
	relationErr := mysql.SelectAllByCond(consts.TableProjectRelation, db.Cond{
		consts.TcOrgId:        orgId,
		consts.TcIsDelete:     consts.AppIsNoDelete,
		consts.TcRelationType: consts.ProjectRelationTypeOwner,
		consts.TcProjectId:    id,
	}, relationPos)
	if relationErr != nil {
		log.Error(relationErr)
		return *projectBo, errs.MysqlOperateError
	}
	for _, relation := range *relationPos {
		projectBo.OwnerIds = append(projectBo.OwnerIds, relation.RelationId)
	}
	return *projectBo, nil
}

func UpdateProject(orgId int64, currentUserId int64, upd mysql.Upd, input *bo.UpdateProjectBo, appId int64) errs.SystemErrorInfo {
	var resourceId int64
	if input.ResourcePath != nil && input.ResourceType != nil {
		respVo := resourcefacade.GetIdByPath(
			resourcevo.GetIdByPathReqVo{
				OrgId:        orgId,
				ResourceType: *input.ResourceType,
				ResourcePath: *input.ResourcePath,
			})
		if !respVo.Failure() {
			resourceId = respVo.ResourceId
		}
	}

	err := mysql.TransX(func(tx sqlbuilder.Tx) error {
		//插入资源
		errSys := updateSource(orgId, currentUserId, input, resourceId, upd)
		if errSys != nil {
			log.Error(errSys)
			return errSys
		}

		//更新项目
		projectError := updateProject(upd, tx, currentUserId, input)
		if projectError != nil {
			log.Error(projectError)
			return projectError
		}
		if newName, ok := upd[consts.TcName]; ok {
			//如果更新项目名称
			if appId > 0 {
				lcResp := appfacade.UpdateLessCodeApp(&appvo.UpdateLessCodeAppReq{
					AppId:  appId,
					OrgId:  orgId,
					UserId: currentUserId,
					Name:   fmt.Sprintf("%v", newName),
				})
				if lcResp.Failure() {
					log.Error(lcResp.Error())
					return lcResp.Error()
				}
			}
		}

		return nil
	})
	if err != nil {
		log.Error(err)
		return errs.BuildSystemErrorInfo(errs.ProjectDomainError, err)
	}
	return nil
}

func updateProject(upd mysql.Upd, tx sqlbuilder.Tx, currentUserId int64, input *bo.UpdateProjectBo) errs.SystemErrorInfo {
	if len(upd) > 0 {
		//更新项目
		upd[consts.TcUpdator] = currentUserId
		upd[consts.TcUpdateTime] = time.Now()
		_, updateProjectErr := mysql.TransUpdateSmartWithCond(tx, consts.TableProject, db.Cond{
			consts.TcId: input.ID,
		}, upd)
		if updateProjectErr != nil {
			return errs.BuildSystemErrorInfo(errs.MysqlOperateError, updateProjectErr)
		}
	}
	return nil
}

func updateSource(orgId, currentUserId int64, input *bo.UpdateProjectBo, resourceId int64, upd mysql.Upd) errs.SystemErrorInfo {
	if util.FieldInUpdate(input.UpdateFields, "resourcePath") && util.FieldInUpdate(input.UpdateFields, "resourceType") {
		if input.ResourcePath != nil && input.ResourceType != nil {
			if resourceId != 0 {
				upd[consts.TcResourceId] = resourceId
			} else {
				fileName := util.ParseFileName(*input.ResourcePath)
				suffix := util.ParseFileSuffix(fileName)
				respVo := resourcefacade.CreateResource(resourcevo.CreateResourceReqVo{
					CreateResourceBo: bo.CreateResourceBo{
						Path:       *input.ResourcePath,
						Name:       fileName,
						Suffix:     suffix,
						OrgId:      orgId,
						OperatorId: currentUserId,
						Type:       *input.ResourceType,
					},
				})
				if respVo.Failure() {
					return respVo.Error()
				}
				upd[consts.TcResourceId] = respVo.ResourceId
			}
		} else {
			upd[consts.TcResourceId] = 0
		}
	}
	return nil
}

func UpdateProjectCondAssembly(input bo.UpdateProjectBo, orgId int64, old, new *map[string]interface{}, originProjectInfo bo.ProjectBo, changeList *[]bo.TrendChangeListBo) (mysql.Upd, errs.SystemErrorInfo) {
	planStartTime := time.Time(originProjectInfo.PlanStartTime)
	planEndTime := time.Time(originProjectInfo.PlanEndTime)
	upd := mysql.Upd{}

	repeatErr := needUpdateVertifyRepeat(input, &upd, orgId, old, new, originProjectInfo, changeList)
	if repeatErr != nil {
		return nil, repeatErr
	}
	// 之后更新走无码，极星不更新优先级，这里不做校验, 2021-07-28
	// deleted annotating code
	needUpdateVertifyValidField(input, &upd, old, new, originProjectInfo, changeList)
	simpleErr := needUpdateSimpleField(input, &upd, old, new, originProjectInfo, changeList)
	if simpleErr != nil {
		return nil, simpleErr
	}
	planTimeErr := needUpdatePlanTime(input, &planStartTime, &planEndTime, &upd, old, new, originProjectInfo, changeList)
	if planTimeErr != nil {
		return nil, planTimeErr
	}

	if util.FieldInUpdate(input.UpdateFields, "preCode") {
		if originProjectInfo.PreCode != "" || input.PreCode == nil {
			return nil, errs.ProjectPreCodeCannotModify
		}

		isPreCodeRight := format.VerifyProjectPreviousCodeFormat(*input.PreCode)
		if !isPreCodeRight {
			log.Error(errs.InvalidProjectPreCodeError)
			return nil, errs.InvalidProjectPreCodeError
		}

		_, err := JudgeRepeatProjectPreCode(input.PreCode, orgId, nil)
		if err != nil {
			log.Error(err)
			return nil, err
		}

		upd[consts.TcPreCode] = input.PreCode
		(*old)["preCode"] = consts.BlankString
		(*new)["preCode"] = input.PreCode
		*changeList = append(*changeList, bo.TrendChangeListBo{
			Field:     "preCode",
			FieldName: consts.ProjectPreCode,
			OldValue:  consts.BlankString,
			NewValue:  *input.PreCode,
		})
	}

	return upd, nil
}

func needUpdateVertifyValidField(input bo.UpdateProjectBo, upd *mysql.Upd, old, new *map[string]interface{}, originProjectInfo bo.ProjectBo, changeList *[]bo.TrendChangeListBo) {
	publicStatus := map[int]string{
		consts.PrivateProject: "私有",
		consts.PublicProject:  "公开",
	}
	if util.FieldInUpdate(input.UpdateFields, "publicStatus") {
		if input.PublicStatus != nil {
			if ok, _ := slice.Contain([]int{consts.PrivateProject, consts.PublicProject}, *input.PublicStatus); ok {
				(*upd)[consts.TcPublicStatus] = input.PublicStatus
				(*old)["publicStatus"] = originProjectInfo.PublicStatus
				(*new)["publicStatus"] = input.PublicStatus
				*changeList = append(*changeList, bo.TrendChangeListBo{
					Field:     "publicStatus",
					FieldName: consts.PublicStatus,
					OldValue:  publicStatus[originProjectInfo.PublicStatus],
					NewValue:  publicStatus[*input.PublicStatus],
				})
			}
		}
	}

	if util.FieldInUpdate(input.UpdateFields, "isFiling") {
		//todo 暂时归档项目不放在更新项目的动态
		if input.IsFiling != nil {
			if ok, _ := slice.Contain([]int{consts.ProjectIsFiling, consts.ProjectIsNotFiling}, *input.IsFiling); ok {
				(*upd)[consts.TcIsFiling] = input.IsFiling
				(*old)["isFiling"] = originProjectInfo.IsFiling
				(*new)["isFiling"] = input.IsFiling
			}
		}
	}
}

func needUpdateVertifyRepeat(input bo.UpdateProjectBo, upd *mysql.Upd, orgId int64, old, new *map[string]interface{}, originProjectInfo bo.ProjectBo, changeList *[]bo.TrendChangeListBo) errs.SystemErrorInfo {
	//判断项目名是否重复
	if util.FieldInUpdate(input.UpdateFields, "name") {
		if input.Name == nil {
			return nil
		} else if strings.Trim(*input.Name, " ") == consts.BlankString {
			return errs.ProjectNameEmpty
		}
		isNameRight := format.VerifyProjectNameFormat(*input.Name)
		if !isNameRight {
			log.Error(errs.InvalidProjectNameError)
			return errs.InvalidProjectNameError
		}

		(*upd)[consts.TcName] = *input.Name
		(*old)["name"] = originProjectInfo.Name
		(*new)["name"] = input.Name
		*changeList = append(*changeList, bo.TrendChangeListBo{
			Field:     "title",
			FieldName: consts.Title,
			OldValue:  originProjectInfo.Name,
			NewValue:  *input.Name,
		})
	}

	return nil
}

func needUpdateSimpleField(input bo.UpdateProjectBo, upd *mysql.Upd, old, new *map[string]interface{}, originProjectInfo bo.ProjectBo, changeList *[]bo.TrendChangeListBo) errs.SystemErrorInfo {
	if util.FieldInUpdate(input.UpdateFields, "remark") {
		if input.Remark != nil {
			//if strs.Len(*input.Remark) > 500 {
			//	return errs.TooLongProjectRemark
			//}
			isRemarkRight := format.VerifyProjectRemarkFormat(*input.Remark)
			if !isRemarkRight {
				log.Error(errs.InvalidProjectRemarkError)
				return errs.InvalidProjectRemarkError
			}
			(*upd)[consts.TcRemark] = input.Remark
		} else {
			(*upd)[consts.TcRemark] = consts.BlankString
		}
		(*old)["remark"] = originProjectInfo.Remark
		(*new)["remark"] = input.Remark
		*changeList = append(*changeList, bo.TrendChangeListBo{
			Field:     "remark",
			FieldName: consts.Remark,
			OldValue:  originProjectInfo.Remark,
			NewValue:  *input.Remark,
		})
	}
	return nil
}

func needUpdatePlanTime(input bo.UpdateProjectBo, planStartTime, planEndTime *time.Time, upd *mysql.Upd, old, new *map[string]interface{}, originProjectInfo bo.ProjectBo, changeList *[]bo.TrendChangeListBo) errs.SystemErrorInfo {
	if util.FieldInUpdate(input.UpdateFields, "planStartTime") {
		if input.PlanStartTime != nil && input.PlanStartTime.IsNotNull() {
			(*upd)[consts.TcPlanStartTime] = date.FormatTime(*input.PlanStartTime)
			*planStartTime = time.Time(*input.PlanStartTime)
		} else {
			(*upd)[consts.TcPlanStartTime] = consts.BlankTime
			*planStartTime = consts.BlankTimeObject
		}
		(*old)["planStartTime"] = originProjectInfo.PlanStartTime
		(*new)["planStartTime"] = (*upd)[consts.TcPlanStartTime]
		*changeList = append(*changeList, bo.TrendChangeListBo{
			Field:     "planStartTime",
			FieldName: consts.PlanStartTime,
			OldValue:  originProjectInfo.PlanStartTime.String(),
			NewValue:  (*upd)[consts.TcPlanStartTime].(string),
		})
	}

	if util.FieldInUpdate(input.UpdateFields, "planEndTime") {
		if input.PlanEndTime != nil && input.PlanEndTime.IsNotNull() {
			(*upd)[consts.TcPlanEndTime] = date.FormatTime(*input.PlanEndTime)
			*planEndTime = time.Time(*input.PlanEndTime)
		} else {
			(*upd)[consts.TcPlanEndTime] = consts.BlankTime
			*planEndTime = consts.BlankTimeObject
		}
		(*old)["planEndTime"] = originProjectInfo.PlanEndTime
		(*new)["planEndTime"] = (*upd)[consts.TcPlanEndTime]
		*changeList = append(*changeList, bo.TrendChangeListBo{
			Field:     "planEndTime",
			FieldName: consts.PlanEndTime,
			OldValue:  originProjectInfo.PlanEndTime.String(),
			NewValue:  (*upd)[consts.TcPlanEndTime].(string),
		})
	}

	if (*planEndTime).After(consts.BlankTimeObject) && planStartTime.After(*planEndTime) {
		return errs.BuildSystemErrorInfo(errs.CreateProjectTimeError)
	}

	return nil
}

func JudgeProjectIsExist(orgId, id int64) bool {
	_, err := LoadProjectAuthBo(orgId, id)
	if err != nil {
		return false
	}

	return true
}

func UpdateProjectStatus(projectBo bo.ProjectBo, nextStatusId int64) errs.SystemErrorInfo {
	orgId := projectBo.OrgId
	projectId := projectBo.Id

	if projectBo.Status == nextStatusId {
		return nil
		//log.Error("更新项目状态-要更新的状态和当前状态一样")
		//return errs.BuildSystemErrorInfo(errs.ProjectStatusUpdateError)
	}

	//验证状态有效性
	if ok, _ := slice.Contain([]int64{consts.StatusRunning.ID, consts.StatusComplete.ID}, nextStatusId); !ok {
		return errs.ProcessStatusNotExist
	}

	_, err2 := dao.UpdateProjectByOrg(projectId, orgId, mysql.Upd{
		consts.TcStatus: nextStatusId,
	})
	if err2 != nil {
		log.Error(err2)
		return errs.BuildSystemErrorInfo(errs.IterationStatusUpdateError)
	}

	return nil
}

// 将前端传过来的日历同步配置，转换为一个配置值，存放在数据库中。
func TransferSyncOutCalendarStatusIntoOne(statusList []*int) int {
	result := consts.IsNotSyncOutCalendar
	if len(statusList) < 1 {
		return result
	}
	sum := int(0)
	for _, item := range statusList {
		sum += *item
	}
	return sum
}

// 日历同步的状态转换。因为数据库中存储的是一个聚合的状态值，需要将其转换为单独的状态值集合。
func TransferSyncOutCalendarStatus(syncFlag int) []*int {
	result := make([]*int, 0)
	status1 := consts.IsSyncOutCalendarForOwner
	status2 := consts.IsSyncOutCalendarForFollower
	status3 := consts.IsSyncOutCalendarForSubCalendar
	// 兼容旧值：1，2
	if syncFlag == consts.IsSyncOutCalendar {
		result = append(result, &status1, &status2, &status3)
		return result
	} else if syncFlag == consts.IsNotSyncOutCalendar {
		return result
	} else {
		if syncFlag&consts.IsSyncOutCalendarForOwner == consts.IsSyncOutCalendarForOwner {
			result = append(result, &status1)
		}
		if syncFlag&consts.IsSyncOutCalendarForFollower == consts.IsSyncOutCalendarForFollower {
			result = append(result, &status2)
		}
		if syncFlag&consts.IsSyncOutCalendarForSubCalendar == consts.IsSyncOutCalendarForSubCalendar {
			result = append(result, &status3)
		}
		return result
	}
}

// CheckProFsChatSetIsOpen 检查群聊（主动创建的群聊）是否开启
// 1 开启；2没有关联群聊；群聊调整后，该方法是检查是否有关联的群聊。
func CheckProFsChatSetIsOpen(orgId, projectId int64) (int, errs.SystemErrorInfo) {
	mainChatId, err := GetProjectMainChatId(orgId, projectId)
	if err != nil {
		log.Errorf("[CheckProFsChatSetIsOpen] projectId: %d, GetMainChatIdByProjectId err: %v", projectId, err)
		return 2, err
	}
	if mainChatId != "" {
		return 1, nil
	}

	return 2, nil
}

// 我参与的所有项目（需排除已被删除的项目）（因ppm_pro_project_relation未删除已删除的项目的关联数据）
func GetAllMyProjectIdsWithDeptIds(orgId, userId int64, deptIds []int64, needFilterArchive bool) ([]int64, errs.SystemErrorInfo) {
	conn, err := mysql.GetConnect()
	if err != nil {
		return nil, errs.MysqlOperateError
	}
	conn.SetLogging(true)
	// 加入部门0，代表整个组织
	allDeptIds := []int64{0}
	allDeptIds = append(allDeptIds, deptIds...)
	union := db.Or(db.Cond{
		"r." + consts.TcRelationId:   userId,
		"r." + consts.TcRelationType: db.In([]int64{consts.ProjectRelationTypeOwner, consts.ProjectRelationTypeParticipant}),
	}, db.Cond{
		"r." + consts.TcRelationId:   db.In(allDeptIds),
		"r." + consts.TcRelationType: consts.ProjectRelationTypeDepartmentParticipant,
	})

	infos := []po.PpmProProjectRelation{}
	cond := db.Cond{
		"r." + consts.TcProjectId: db.Raw("p." + consts.TcId),
		"r." + consts.TcOrgId:     orgId,
		"r." + consts.TcIsDelete:  consts.AppIsNoDelete,
		"p." + consts.TcIsDelete:  consts.AppIsNoDelete,
	}
	if needFilterArchive {
		cond["p."+consts.TcIsFiling] = consts.ProjectIsNotFiling
	}
	q := conn.Select(db.Raw("r.*")).From(consts.TableProjectRelation+" r", consts.TableProject+" p").Where(cond).And(union)
	log.Infof("[GetAllMyProjectIdsWithDeptIds] %v", q.String())

	err = q.All(&infos)
	if err != nil {
		log.Error(err)
		return nil, errs.MysqlOperateError
	}
	projectIds := []int64{0} //默认把未归属加进来
	for _, info := range infos {
		projectIds = append(projectIds, info.ProjectId)
	}
	return slice.SliceUniqueInt64(projectIds), nil
}

// 过滤掉归档的project和空项目
func GetAllNotFillingAnEmptyProjectIds(orgId int64, isGetAllType bool) ([]int64, errs.SystemErrorInfo) {
	infos := []po.PpmProProject{}
	cond := db.Cond{
		consts.TcIsDelete: consts.AppIsNoDelete,
		consts.TcOrgId:    orgId,
		consts.TcIsFiling: consts.ProjectIsNotFiling,
	}
	if !isGetAllType {
		cond[consts.TcProjectTypeId] = db.NotEq(consts.ProjectTypeEmpty)
	}
	err := mysql.SelectAllByCond(consts.TableProject, cond, &infos)

	if err != nil {
		log.Error(err)
		return nil, errs.MysqlOperateError
	}
	allProjectIds := make([]int64, 0)
	for _, info := range infos {
		allProjectIds = append(allProjectIds, info.Id)
	}
	allProjectIds = append(allProjectIds, 0)

	return allProjectIds, nil
}

func GetAllProjectIdsInterface(orgId int64) ([]interface{}, errs.SystemErrorInfo) {
	infos := []po.PpmProProject{}
	err := mysql.SelectAllByCond(consts.TableProject, db.Cond{
		consts.TcIsDelete: consts.AppIsNoDelete,
		consts.TcOrgId:    orgId,
	}, &infos)
	if err != nil {
		log.Error(err)
		return nil, errs.MysqlOperateError
	}
	allProjectIds := make([]interface{}, 0)
	for _, info := range infos {
		allProjectIds = append(allProjectIds, info.Id)
	}
	allProjectIds = append(allProjectIds, 0)

	return allProjectIds, nil
}

func convertStringSliceToInt64(ss []string) ([]int64, error) {
	is := make([]int64, 0, len(ss))
	for _, s := range ss {
		if i, err := strconv.ParseInt(s, 10, 64); err == nil {
			is = append(is, i)
		}
	}
	return is, nil
}

func convertPbStringListToSliceInt64(l *structpb.ListValue) ([]int64, error) {
	is := make([]int64, 0, len(l.Values))
	for _, s := range l.Values {
		if i, err := strconv.ParseInt(s.GetStringValue(), 10, 64); err == nil {
			is = append(is, i)
		}
	}
	return is, nil
}

func parseMemberStr(member string) (string, int64, bool) {
	ss := strings.Split(member, "_")
	if len(ss) != 2 {
		return "", 0, false
	}
	id, err := strconv.ParseInt(ss[1], 0, 64)
	if err != nil {
		return "", 0, false
	}
	return ss[0], id, true
}

func GetUserCollaboratorRoleIdsByAppId(orgId, appId, userId int64, deptIds []int64, datas []map[string]interface{}) ([]int64, errs.SystemErrorInfo) {
	tableColumnsResp := tablefacade.ReadTableSchemasByAppId(projectvo.GetTableSchemasByAppIdReq{
		OrgId:  orgId,
		UserId: userId,
		Input: &tableV1.ReadTableSchemasByAppIdRequest{
			AppId: appId,
		},
	})
	if tableColumnsResp.Failure() {
		return nil, tableColumnsResp.Error()
	}
	if len(tableColumnsResp.Data.Tables) == 0 {
		return nil, errs.TableNotExist
	}
	return getUserCollaboratorRoleIds(orgId, appId, userId, deptIds, datas, tableColumnsResp.Data.Tables[0])
}

func GetUserCollaboratorRoleIdsByTableId(orgId, appId, tableId, userId int64, deptIds []int64, datas []map[string]interface{}) ([]int64, errs.SystemErrorInfo) {
	tableColumnsResp := tablefacade.GetTableColumns(projectvo.GetTableColumnsReq{
		OrgId:  orgId,
		UserId: userId,
		Input: &tableV1.ReadTableSchemasRequest{
			TableIds: []int64{tableId},
		},
	})
	if tableColumnsResp.Failure() {
		return nil, tableColumnsResp.Error()
	}
	if len(tableColumnsResp.Data.Tables) == 0 || tableColumnsResp.Data.Tables[0].TableId != tableId {
		return nil, errs.TableNotExist
	}
	return getUserCollaboratorRoleIds(orgId, appId, userId, deptIds, datas, tableColumnsResp.Data.Tables[0])

}

// 实现无码APP服务的GetUserCollaboratorRoleIds接口，在极星已经持有任务数据时使用，避免再调用无码重复拉任务数据
func getUserCollaboratorRoleIds(orgId, appId, userId int64, deptIds []int64, datas []map[string]interface{}, tableSchema *projectvo.TableColumnsTable) ([]int64, errs.SystemErrorInfo) {
	roleIds := make([]int64, 0)
	if len(datas) == 0 {
		return roleIds, nil
	}

	collaboratorFields := make(map[string][]int64)
	for _, column := range tableSchema.Columns {
		fieldData := column.Field
		if fieldData.Type == tableV1.ColumnType_dept.String() ||
			fieldData.Type == tableV1.ColumnType_member.String() ||
			fieldData.Type == tableV1.ColumnType_workHour.String() {
			//propsMap := fieldData.Props
			//props := propsMap["field"].(map[string]interface{})
			props := fieldData.Props
			if cRoles, ok := props["collaboratorRoles"]; ok {
				// collaboratorRoles, err := convertPbStringListToSliceInt64(cRoles.GetListValue()) // cRoles.GetListValue()
				collaboratorRoles := make([]int64, 0)
				tmpJson := json.ToJsonIgnoreError(cRoles)
				err := json.FromJson(tmpJson, &collaboratorRoles)
				if err != nil {
					return nil, errs.JSONConvertError
				}
				collaboratorFields[column.Name] = collaboratorRoles
			}
		}
	}
	if len(collaboratorFields) == 0 {
		return roleIds, nil
	}

	// 所有部门ID，加上0
	allDeptIds := make(map[int64]struct{})
	allDeptIds[0] = struct{}{}
	for _, id := range deptIds {
		allDeptIds[id] = struct{}{}
	}

	// 获取所有权限组
	roleListResp := permissionfacade.GetAppRoleList(orgId, appId)
	if roleListResp.Failure() {
		return nil, roleListResp.Error()
	}
	var editRoleId, finalEditRoleId int64
	for _, role := range roleListResp.Data {
		if role.LangCode == consts.GroupLandCodeProjectMember ||
			role.LangCode == consts.GroupLandCodeEdit {
			// -1是默认的负责人，关注人，确认人等 目前阶段默认为编辑者
			editRoleId = role.Id
			break
		}
	}
	// 把-1替换成编辑者角色
	if editRoleId != 0 {
		finalEditRoleId = editRoleId
		for _, rs := range collaboratorFields {
			for i, r := range rs {
				if r == -1 {
					rs[i] = finalEditRoleId
				}
			}
		}
	}

	// 扫所有数据，找到所有字段上包含自己的协作者角色
	for _, data := range datas {
		for k, v := range data {
			// 协作人字段
			if rIds, ok := collaboratorFields[k]; ok {
				switch vv := v.(type) {
				case []string:
					// 协作人判断
					for _, _vv := range vv {
						if tp, id, ok := parseMemberStr(_vv); ok {
							if tp == "U" && id == userId {
								roleIds = append(roleIds, rIds...)
							} else if tp == "D" {
								if _, ok := allDeptIds[id]; ok {
									roleIds = append(roleIds, rIds...)
								}
							}
						}
					}
				}
			}
		}
	}

	return roleIds, nil
}

//func GetOneProjectInfo(orgId, id int64) (*bo.Project, errs.SystemErrorInfo) {
//
//	projectBo := &bo.Project{}
//
//	projet, err := dao.SelectOneProject(db.Cond{
//		consts.TcId:       id,
//		consts.TcOrgId:    orgId,
//		consts.TcIsDelete: consts.AppIsNoDelete,
//	})
//	if err != nil {
//		return nil, errs.BuildSystemErrorInfo(errs.ProjectNotExist)
//	}
//
//	copyer.Copy(projet, projectBo)
//	return projectBo, nil
//}

func GetProjectAppIdsByProjectIds(orgId int64, projectIds []int64) ([]int64, errs.SystemErrorInfo) {
	proAppIds := make([]int64, 0, len(projectIds))
	if len(projectIds) < 1 {
		return proAppIds, nil
	}
	proList, _, err := GetProjectList(0, db.Cond{
		consts.TcOrgId: orgId,
		consts.TcId:    db.In(projectIds),
	}, nil, nil, 1000, 1)
	if err != nil {
		log.Errorf("[GetTableListByIds] err: %v, orgId: %d", err, orgId)
		return nil, err
	}
	for _, pro := range proList {
		proAppIds = append(proAppIds, pro.AppId)
	}

	return proAppIds, nil
}

func GetBasicFields(projectTypeId int) ([]string, errs.SystemErrorInfo) {
	if projectTypeId == consts.ProjectTypeEmpty {
		return []string{}, nil
	}
	return []string{"_field_priority"}, nil
}

// GetOrgSummaryAppId 获取某个组织的汇总表 appId
func GetOrgSummaryAppId(orgId int64) (int64, errs.SystemErrorInfo) {
	orgResp := orgfacade.GetOrgBoListByPage(orgvo.GetOrgIdListByPageReqVo{
		Page: 1,
		Size: 10000, // 不要超过1w
		Input: orgvo.GetOrgIdListByPageReqVoData{
			OrgIds: []int64{orgId},
		},
	})
	if orgResp.Failure() {
		log.Errorf("[GetOrgSummaryAppId] err: %v", orgResp.Error())
		return 0, orgResp.Error()
	}
	if len(orgResp.Data.List) < 1 {
		return 0, nil
	}
	org := orgResp.Data.List[0]
	orgRemarkObj := &orgvo.OrgRemarkConfigType{}
	if len(org.Remark) > 0 {
		oriErr := json.FromJson(org.Remark, orgRemarkObj)
		if oriErr != nil {
			log.Errorf("[GetOrgSummaryAppId] 组织 remark 反序列化异常，组织id:%d,原因:%v", org.Id, oriErr)
			return 0, errs.BuildSystemErrorInfo(errs.JSONConvertError, oriErr)
		}
	}

	return orgRemarkObj.OrgSummaryTableAppId, nil
}

func GetProjectIdByAppId(orgId, appId int64) (int64, errs.SystemErrorInfo) {
	var projectPo po.PpmProProject
	err := mysql.SelectOneByCond(consts.TableProject, db.Cond{
		consts.TcOrgId:    orgId,
		consts.TcIsDelete: consts.AppIsNoDelete,
		consts.TcAppId:    appId,
	}, &projectPo)
	if err != nil {
		log.Errorf("[GetProjectIdByAppId] err:%v, orgId:%v, appId:%v", err, orgId, appId)
		return 0, errs.MysqlOperateError
	}
	return projectPo.Id, nil
}

func GetProjectIdsByAppIds(orgId int64, appIds []int64) ([]int64, errs.SystemErrorInfo) {
	projectIds := make([]int64, 0, len(appIds))
	projectPos := []po.PpmProProject{}
	err := mysql.SelectAllByCond(consts.TableProject, db.Cond{
		consts.TcOrgId:    orgId,
		consts.TcIsDelete: consts.AppIsNoDelete,
		consts.TcAppId:    db.In(appIds),
	}, &projectPos)
	if err != nil {
		log.Errorf("[GetProjectIdsByAppIds] err:%v, orgId:%v, appIds:%v", err, orgId, appIds)
		return projectIds, errs.MysqlOperateError
	}
	for _, pro := range projectPos {
		projectIds = append(projectIds, pro.Id)
	}
	return projectIds, nil
}

func CreateExternalApp(orgId, userId int64, name, icon, linkUrl string) *permissionvo.CreateLessCodeAppResp {
	appType := consts.LcAppTypeForFolder
	return appfacade.CreateLessCodeApp(&permissionvo.CreateLessCodeAppReq{
		AppType:      &appType,
		OrgId:        &orgId,
		UserId:       &userId,
		Name:         &name,
		AuthType:     2,
		Icon:         icon,
		ExternalApp:  1,
		LinkUrl:      linkUrl,
		AddAllMember: true,
	})
}
