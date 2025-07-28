package resourcesvc

import (
	"fmt"
	"strconv"
	"strings"

	sdk_const "gitea.bjx.cloud/allstar/platform-sdk/consts"

	tableV1 "gitea.bjx.cloud/LessCode/interface/golang/table/v1"
	"gitea.bjx.cloud/allstar/dingtalk-sdk-golang"
	vo2 "gitea.bjx.cloud/allstar/feishu-sdk-golang/core/model/vo"
	platform_sdk "gitea.bjx.cloud/allstar/platform-sdk"
	sdkVo "gitea.bjx.cloud/allstar/platform-sdk/vo"
	"gitea.bjx.cloud/allstar/polaris-backend/facade/idfacade"
	"github.com/star-table/startable-server/app/facade/orgfacade"
	"github.com/star-table/startable-server/app/facade/tablefacade"
	"github.com/star-table/startable-server/common/core/consts"
	"github.com/star-table/startable-server/common/core/errs"
	"github.com/star-table/startable-server/common/core/logger"
	"github.com/star-table/startable-server/common/core/util"
	"github.com/star-table/startable-server/common/core/util/copyer"
	"github.com/star-table/startable-server/common/core/util/json"
	"github.com/star-table/startable-server/common/core/util/slice"
	"github.com/star-table/startable-server/common/core/util/str"
	"github.com/star-table/startable-server/common/core/util/uuid"
	"github.com/star-table/startable-server/common/extra/feishu"
	"github.com/star-table/startable-server/common/library/cache"
	"github.com/star-table/startable-server/common/library/db/mysql"
	"github.com/star-table/startable-server/common/model/bo"
	"github.com/star-table/startable-server/common/model/vo"
	"github.com/star-table/startable-server/common/model/vo/orgvo"
	"github.com/star-table/startable-server/common/model/vo/resourcevo"
	"upper.io/db.v3"
	"upper.io/db.v3/lib/sqlbuilder"
)

var log = logger.GetDefaultLogger()

func CreateResource(createResourceBo bo.CreateResourceBo, tx ...sqlbuilder.Tx) (int64, errs.SystemErrorInfo) {
	return domain.CreateResource(createResourceBo, tx...)
}

func UpdateResourceInfo(input bo.UpdateResourceInfoBo) (*resourcevo.UpdateResourceData, errs.SystemErrorInfo) {
	orgId := input.OrgId
	resourceId := input.ResourceId
	projectId := input.ProjectId
	updateFields := input.UpdateFields
	resp := &resourcevo.UpdateResourceData{}
	if updateFields == nil || len(updateFields) == 0 {
		return nil, errs.UpdateFiledIsEmpty
	}
	err := domain.CheckResourceIds([]int64{resourceId}, projectId, orgId)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	bos, err := domain.GetResourceByIds([]int64{resourceId})
	if err != nil {
		log.Error(err)
		return nil, err
	}
	oldBo := bos[0]
	resp.OldBo = append(resp.OldBo, oldBo)
	newPo, err := domain.UpdateResourceInfo(resourceId, input)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	newBo := &bo.ResourceBo{}
	err1 := copyer.Copy(newPo, newBo)
	if err1 != nil {
		log.Error(err1)
		return nil, errs.BuildSystemErrorInfo(errs.ObjectCopyError, err1)
	}
	resp.NewBo = append(resp.NewBo, *newBo)
	return resp, nil
}

func UpdateResourceFolder(input bo.UpdateResourceFolderBo) (*resourcevo.UpdateResourceData, errs.SystemErrorInfo) {
	orgId := input.OrgId
	userId := input.UserId
	resourceIds := input.ResourceIds
	projectId := input.ProjectId
	currentFolderId := input.CurrentFolderId
	targetFolderId := input.TargetFolderID
	resp := &resourcevo.UpdateResourceData{}
	err := domain.CheckFolderIds([]int64{currentFolderId, targetFolderId}, projectId, orgId)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	err = domain.UpdateResourceFolderId(resourceIds, currentFolderId, targetFolderId, userId, orgId)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	bos, err := domain.GetResourceByIds(resourceIds)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	blank := ""
	if currentFolderId != 0 {
		currentBos, err := domain.GetFolderById([]int64{currentFolderId})
		if err != nil {
			log.Error(err)
			return nil, err
		}
		resp.CurrentFolderName = &currentBos[0].Name
	} else {
		resp.CurrentFolderName = &blank
	}
	if targetFolderId != 0 {
		targetBos, err := domain.GetFolderById([]int64{targetFolderId})
		if err != nil {
			log.Error(err)
			return nil, err
		}
		resp.TargetFolderName = &targetBos[0].Name
	} else {
		resp.TargetFolderName = &blank
	}
	resp.OldBo = bos
	return resp, nil
}

func CompleteDeleteResource(orgId, userId int64, resourceIds []int64) errs.SystemErrorInfo {
	// 引用文件不删除资源本身，只把资源关联关系删除
	isDelete := consts.AppIsNoDelete
	resourceBos, _, errSys := domain.GetResourceBoList(0, 0, resourcevo.GetResourceBoListCond{
		OrgId:       orgId,
		ResourceIds: &resourceIds,
		IsDelete:    &isDelete,
	})
	if errSys != nil {
		log.Errorf("[CompleteDeleteResource] GetResourceBoList err:%v, orgId:%d, resourceIds:%v", errSys, orgId, resourceIds)
		return errSys
	}
	needDeleteIds := []int64{}
	for _, resource := range *resourceBos {
		if resource.SourceType != consts.OssPolicyTypeProjectResource {
			needDeleteIds = append(needDeleteIds, resource.Id)
		}
	}
	if len(needDeleteIds) > 0 {
		_, err := mysql.UpdateSmartWithCond(consts.TableResource, db.Cond{
			consts.TcIsDelete: consts.AppIsNoDelete,
			consts.TcOrgId:    orgId,
			consts.TcId:       db.In(needDeleteIds),
		}, mysql.Upd{
			consts.TcIsDelete: consts.AppIsDeleted,
			consts.TcUpdator:  userId,
		})
		if err != nil {
			log.Error(err)
			return errs.MysqlOperateError
		}
	}

	cacheErr := domain.ClearCacheResourceSize(orgId)
	if cacheErr != nil {
		log.Error(cacheErr)
		return cacheErr
	}

	return nil
}

// 目前项目 文件管理中删除文件会使用到
func DeleteResource(deleteBo bo.DeleteResourceBo) (*resourcevo.UpdateResourceData, errs.SystemErrorInfo) {
	orgId := deleteBo.OrgId
	userId := deleteBo.UserId
	resourceIds := deleteBo.ResourceIds
	folderId := deleteBo.FolderId
	//仅文件做文件夹校验
	//if folderId != nil {
	//	err := domain.CheckRelation(resourceIds, *folderId, orgId)
	//	if err != nil {
	//		log.Error(err)
	//		return nil, err
	//	}
	//}
	resp := &resourcevo.UpdateResourceData{}
	bos, err := domain.GetResourceByIds(resourceIds)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	transErr := mysql.TransX(func(tx sqlbuilder.Tx) error {
		return domain.DeleteResource(resourceIds, folderId, orgId, deleteBo.AppId, userId,
			deleteBo.ProjectId, deleteBo.IssueId, deleteBo.RecycleVersionId, tx)
	})
	if transErr != nil {
		log.Errorf("[DeleteResource] DeleteResource err:%v", transErr)
		return nil, errs.MysqlOperateError
	}
	resp.OldBo = bos
	return resp, nil
}

func GetResource(input bo.GetResourceBo) (*vo.ResourceList, errs.SystemErrorInfo) {
	folderId := input.FolderId
	orgId := input.OrgId
	projectId := input.ProjectId
	cond := db.Cond{
		consts.TcIsDelete: consts.AppIsNoDelete,
		//consts.TcProjectId: projectId,
		consts.TcOrgId: orgId,
		//新增获取的文件类型 2019/12/27
		consts.TcSourceType: db.In(input.SourceType),
	}
	if input.FileType != nil {
		cond[consts.TcFileType] = *input.FileType
	}
	if input.KeyWord != nil {
		cond[consts.TcName] = db.Like("%" + *input.KeyWord + "%")
	}
	pageBo := bo.PageBo{Page: input.Page, Size: input.Size, Order: "id desc"}

	// 查询文件夹的资源
	if folderId != nil {
		err := domain.CheckFolderIds([]int64{*folderId}, projectId, orgId)
		if err != nil {
			log.Error(err)
			return nil, err
		}
		resourceIds, err := domain.GetResourceIdsByFolderId(*folderId, orgId)
		if err != nil {
			log.Error(err)
			return nil, err
		}
		cond[consts.TcId] = db.In(resourceIds)
	}
	// 之前在文件管理上传项目文件，没有在表ppm_res_resource表记录project_id，需要从资源关联表中找
	cond[consts.TcId+" "] = db.In(db.Raw("select distinct resource_id from ppm_res_resource_relation where project_id = ? and is_delete = 2 and source_type in ?", projectId, input.SourceType))

	resourceBos, total, err := domain.GetResourceBoListByPage(cond, nil, pageBo)
	if err != nil {
		log.Error(err)
		return nil, err
	}

	resourceVos := &[]*vo.Resource{}
	copyErr := copyer.Copy(resourceBos, resourceVos)
	if copyErr != nil {
		log.Error(copyErr)
		return nil, errs.BuildSystemErrorInfo(errs.ObjectCopyError, copyErr)
	}

	creatorIds := make([]int64, 0)
	for _, value := range *resourceVos {
		creatorIds = append(creatorIds, value.Creator)
		value.PathCompressed = util.GetCompressedPath(value.Host+value.Path, value.Type)
	}
	ownerMap, err := domain.GetBaseUserInfoMap(orgId, creatorIds)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	for i, resourceVo := range *resourceVos {
		if ownerInfoInterface, ok := ownerMap[resourceVo.Creator]; ok {
			ownerInfo := ownerInfoInterface.(bo.BaseUserInfoBo)
			(*resourceVos)[i].CreatorName = ownerInfo.Name
		} else {
			log.Errorf("用户 %d 信息不存在，组织id %d", resourceVo.Creator, orgId)
		}
	}
	return &vo.ResourceList{
		List:  *resourceVos,
		Total: total,
	}, nil
}

func GetResourceInfo(input bo.GetResourceInfoBo) (*vo.Resource, errs.SystemErrorInfo) {
	orgId := input.OrgId
	//projectId := input.AppId
	resourceId := input.ResourceId
	cond := db.Cond{
		consts.TcIsDelete: consts.AppIsNoDelete,
		//consts.TcProjectId: projectId,
		consts.TcOrgId: orgId,
		//新增获取的文件类型 2019/12/27
		consts.TcSourceType: db.In(input.SourceTypes),
	}

	cond[consts.TcId] = db.In([]int64{resourceId})
	// cond[consts.TcId+" "] = db.In(db.Raw("select distinct id from ppm_res_resource_relation where project_id = ? and resource_id = ? and is_delete = 2 and source_type = ?", projectId, resourceId, input.SourceType))

	pageBo := bo.PageBo{Page: 1, Size: 20, Order: "id desc"}

	resourceBos, _, err := domain.GetResourceBoListByPage(cond, nil, pageBo)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	resourceVos := &[]*vo.Resource{}
	copyErr := copyer.Copy(resourceBos, resourceVos)
	if copyErr != nil {
		log.Error(copyErr)
		return nil, errs.BuildSystemErrorInfo(errs.ObjectCopyError, copyErr)
	}
	creatorIds := make([]int64, 0)
	for _, value := range *resourceVos {
		creatorIds = append(creatorIds, value.Creator)
		value.PathCompressed = util.GetCompressedPath(value.Host+value.Path, value.Type)
	}
	ownerMap, err := domain.GetBaseUserInfoMap(orgId, creatorIds)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	for i, resourceVo := range *resourceVos {
		if ownerInfoInterface, ok := ownerMap[resourceVo.Creator]; ok {
			ownerInfo := ownerInfoInterface.(bo.BaseUserInfoBo)
			(*resourceVos)[i].CreatorName = ownerInfo.Name
		} else {
			log.Errorf("用户 %d 信息不存在，组织id %d", resourceVo.Creator, orgId)
		}
	}

	if len(*resourceVos) == 0 {
		return &vo.Resource{}, nil
	}

	resource := (*resourceVos)[0]
	return &vo.Resource{
		ID: resource.ID,
		// 组织id
		OrgID: resource.OrgID,
		// host
		Host: resource.Host,
		// 路径
		Path: resource.Path,
		// 缩略图路径
		PathCompressed: resource.PathCompressed,
		// 文件名
		Name: resource.Name,
		// 存储类型,1：本地，2：oss,3.钉盘
		Type: resource.Type,
		// 文件大小
		Size: resource.Size,
		// 创建人姓名
		CreatorName: resource.CreatorName,
		// 文件后缀
		Suffix: resource.Suffix,
		// 文件的md5
		Md5: resource.Md5,
		// 文件类型
		FileType: resource.FileType,
		// 创建人
		Creator: resource.Creator,
		// 创建时间
		CreateTime: resource.CreateTime,
		// 更新人
		Updator: resource.Updator,
		// 更新时间
		UpdateTime: resource.UpdateTime,
		// 乐观锁
		Version: resource.Version,
		// 是否删除,1是,2否
		IsDelete: resource.IsDelete,
	}, nil
}

//func InsertResource(tx sqlbuilder.Tx, resourcePath string, orgId int64, currentUserId int64, resourceType int, fileName string) (int64, errs.SystemErrorInfo) {
//	return domain.InsertResource(tx, resourcePath, orgId, currentUserId, resourceType, fileName)
//}

// 获取资源信息
func GetResourceById(resourceIds []int64) ([]bo.ResourceBo, errs.SystemErrorInfo) {
	return domain.GetResourceByIds(resourceIds)
}

func GetIdByPath(orgId int64, resourcePath string, resourceType int) (int64, errs.SystemErrorInfo) {
	return domain.GetIdByPath(orgId, resourcePath, resourceType)
}

func RecoverResource(orgId, appId, userId, projectId, relationId int64, recycleVersionId int64, issueIds []int64, sourceChannel string) (*bo.ResourceBo, errs.SystemErrorInfo) {
	// 恢复前先查一下原先资源所在文件夹是否也被删除了
	isDelete := domain.CheckResourceFolderIsDelete(orgId, projectId, relationId, int(recycleVersionId))
	if isDelete {
		return nil, errs.RecoverDocumentFailedWithNoFolder
	}

	res := &bo.ResourceBo{}
	if relationId != 0 {
		info, err := domain.GetResourceInfo(orgId, relationId)
		if err != nil {
			log.Error(err)
			return nil, err
		}
		res = info
	}

	transErr := mysql.TransX(func(tx sqlbuilder.Tx) error {
		//恢复文件夹文件关联
		if relationId != 0 {
			_, err := mysql.TransUpdateSmartWithCond(tx, consts.TableFolderResource, db.Cond{
				consts.TcIsDelete: consts.AppIsDeleted,
				consts.TcVersion:  recycleVersionId,
				//consts.TcProjectId:projectId,
				consts.TcOrgId:      orgId,
				consts.TcResourceId: relationId,
			}, mysql.Upd{
				consts.TcIsDelete: consts.AppIsNoDelete,
				consts.TcUpdator:  userId,
			})
			if err != nil {
				log.Error(err)
				return err
			}
		}

		//恢复文件关联关系
		cond := db.Cond{
			consts.TcIsDelete:  consts.AppIsDeleted,
			consts.TcVersion:   recycleVersionId,
			consts.TcProjectId: projectId,
			consts.TcOrgId:     orgId,
		}
		if relationId != 0 {
			cond[consts.TcResourceId] = relationId
		}
		if issueIds != nil && len(issueIds) > 0 {
			cond[consts.TcIssueId] = db.In(issueIds)
		}

		_, err := mysql.TransUpdateSmartWithCond(tx, consts.TableResourceRelation, cond, mysql.Upd{
			consts.TcIsDelete: consts.AppIsNoDelete,
			consts.TcUpdator:  userId,
		})
		if err != nil {
			log.Error(err)
			return err
		}

		//cond[consts.TcIsDelete] = consts.AppIsNoDelete
		//return recoverResourceToLess(orgId, appId, userId, cond, tx)
		//return domain.RecoverLessAttachment(orgId, userId, appId, recycleId, tx)
		return domain.RecoverLessAttachments(orgId, userId, projectId, appId, recycleVersionId, sourceChannel, tx)
	})

	if transErr != nil {
		log.Error(transErr)
		if transErr.Error() == errs.IssueAlreadyBeDeleted.Error() {
			return nil, errs.RecoverAttachmentError
		}
		if transErr.Error() == errs.CanNotRecoverDocuments.Error() {
			return nil, errs.CanNotRecoverDocuments
		}
		return nil, errs.RecoverResourceFailed
	}

	return res, nil
}

func recoverResourceToLess(orgId, appId, userId int64, cond db.Cond, tx sqlbuilder.Tx) error {
	if appId == 0 {
		return nil
	}

	issueIds, resourceIds, err := domain.GetRelationIssueIdsAndResourceIds(cond, tx)
	if err != nil {
		return errs.MysqlOperateError
	}
	if len(issueIds) == 0 {
		return nil
	}

	resp := tablefacade.RecoverAttachment(orgId, userId, &tableV1.RecoverAttachmentRequest{AppId: appId, IssueIds: issueIds, ResourceIds: resourceIds})
	if resp.Failure() {
		log.Errorf("[recoverResourceToLess] RecoverAttachment error:%v", resp.Err)
		return errs.SystemError
	}

	return nil
}

func CacheResourceSize(orgId int64) (int64, errs.SystemErrorInfo) {
	return domain.CacheResourceSize(orgId)
}

func FsDocList(orgId int64, userId int64, page int, size int, searchKey string) (*vo.FsDocumentListResp, errs.SystemErrorInfo) {
	baseOrgInfo, err := orgfacade.GetBaseOrgInfoRelaxed(orgId)
	if err != nil {
		log.Errorf("[FsDocList] orgfacade.GetBaseOrgInfoRelaxed failed, orgId:%d, userId:%d, err: %v", orgId, userId, err)
		return nil, err
	}

	client, errClient := platform_sdk.GetClient(sdk_const.SourceChannelFeishu, baseOrgInfo.OutOrgId)
	if errClient != nil {
		return nil, errs.FeiShuClientTenantError
	}
	accessTokenResp := orgfacade.GetFsAccessToken(orgvo.GetFsAccessTokenReqVo{UserId: userId, OrgId: orgId})
	if accessTokenResp.Failure() {
		log.Errorf("[FsFileList] orgfacade.GetFsAccessToken failed, orgId:%d, userId:%d, err: %v", orgId, userId, accessTokenResp.Error())
		return nil, accessTokenResp.Error()
	}
	fsReq := &sdkVo.DocumentReq{
		SearchKey:       searchKey,
		UserAccessToken: accessTokenResp.AccessToken,
		DocsTypes:       []string{consts.FsDocumentTypeDoc, consts.FsDocumentTypeFile, consts.FsDocumentTypeSheet, consts.FsDocumentTypeMindNote, consts.FsDocumentTypeBiTable},
		Page:            page,
		Size:            size,
	}
	fileList, sdkError := client.GetDocumentList(fsReq)
	if sdkError != nil {
		log.Errorf("[DingFileList]failed, orgId: %d, userId: %d, err:%v", orgId, userId, sdkError)
		return nil, errs.FeiShuOpenApiCallError
	}

	openIds := []string{}

	for _, f := range fileList.List {
		openIds = append(openIds, f.OwnerId)
	}
	openIds = slice.SliceUniqueString(openIds)

	userMap := map[string]string{}
	if len(openIds) > 0 {
		// 批量获取用户信息
		userBatchInfoResp := orgfacade.GetOrgUsersInfoByEmIds(orgvo.GetOrgUsersInfoByEmIdsReq{
			OrgId:         orgId,
			SourceChannel: sdk_const.SourceChannelDingTalk,
			EmpIds:        openIds,
		})
		if userBatchInfoResp.Failure() {
			log.Errorf("[FsDocList] orgfacade.GetOrgUsersInfoByEmIds failed, orgId: %d, userId: %d, err: %v",
				orgId, userId, userBatchInfoResp.Error())
			return nil, userBatchInfoResp.Error()
		}
		for _, u := range userBatchInfoResp.BaseUserInfos {
			userMap[u.OutUserId] = u.Name
		}
	}

	list := make([]*vo.FsDocumentData, fileList.Total)
	for _, file := range fileList.List {
		title := file.Title
		if title == "" {
			title = "未命名文档"
		}
		temp := &vo.FsDocumentData{
			Title:     title,
			DocsType:  file.DocsType,
			DocsToken: file.Token,
			DocsURL:   buildDocumentUrl(file.DocsType, file.Token),
			OwnerID:   file.OwnerId,
		}
		if _, ok := userMap[file.OwnerId]; ok {
			temp.OwnerName = userMap[file.OwnerId]
		}
		list = append(list, temp)
	}

	return &vo.FsDocumentListResp{
		Total: int64(fileList.Total),
		List:  list,
	}, nil
}

func FsDocumentList(orgId int64, userId int64, page int, size int, searchKey string) (*vo.FsDocumentListResp, errs.SystemErrorInfo) {
	baseOrgInfo, err := orgfacade.GetBaseOrgInfoRelaxed(orgId)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	// 检查是否有飞书云文档权限，如果没有则提示。前端调用接口检查过，这里无需再调用。
	//resp := orgfacade.CheckSpecificScope(orgvo.CheckSpecificScopeReqVo{
	//	OrgId:     orgId,
	//	UserId:    userId,
	//	PowerFlag: consts.FsScopeDriveRead,
	//})
	//if resp.Failure() {
	//	return nil, errs.FeiShuScopeNeedApply
	//}
	if baseOrgInfo.OutOrgId == "" {
		return nil, errs.CannotBindChat
	}
	tenant, err1 := feishu.GetTenant(baseOrgInfo.OutOrgId)
	if err1 != nil {
		log.Error(err1)
		return nil, err1
	}

	accessTokenResp := orgfacade.GetFsAccessToken(orgvo.GetFsAccessTokenReqVo{UserId: userId, OrgId: orgId})
	if accessTokenResp.Failure() {
		log.Error(accessTokenResp.Error())
		return nil, accessTokenResp.Error()
	}

	req := vo2.SearchDocsReqVo{
		//目前只取文档，表单，文件这三种，因为只有这三种可推断出飞书的url，哎~~~~~~~
		DocsTypes: &[]string{consts.FsDocumentTypeDoc, consts.FsDocumentTypeDocx, consts.FsDocumentTypeFile, consts.FsDocumentTypeSheet, consts.FsDocumentTypeMindNote, consts.FsDocumentTypeBiTable},
	}
	req.SearchKey = &searchKey
	//documentResp, documentErr := tenant.SearchDocs(accessTokenResp.AccessToken, req)
	//if documentErr != nil {
	//	log.Error(documentErr)
	//	return nil, errs.FeiShuOpenApiCallError
	//}
	//
	//if documentResp.Code != 0 {
	//	log.Error(json.ToJsonIgnoreError(documentResp))
	//	return nil, errs.FeiShuOpenApiCallError
	//}

	list := []vo2.DocEntities{}
	total := int64(0)
	if page <= 0 || size <= 0 || size > 50 {
		//表示查询全部(默认单次最大50条)
		count := 50
		req.Count = &count
		i := 1
		for {
			offset := (i - 1) * count
			req.Offset = &offset
			documentResp, documentErr := tenant.SearchDocs(accessTokenResp.AccessToken, req)
			if documentErr != nil {
				log.Error(documentErr)
				return nil, errs.FeiShuOpenApiCallError
			}

			if documentResp.Code != 0 {
				log.Error(json.ToJsonIgnoreError(documentResp))
				return nil, errs.FeiShuOpenApiCallError
			}

			total = documentResp.Data.Total
			list = append(list, documentResp.Data.DocsEntities...)

			if !documentResp.Data.HasMore {
				break
			}

			i++
		}
	} else {
		req.Count = &size
		offset := (page - 1) * size
		req.Offset = &offset
		documentResp, documentErr := tenant.SearchDocs(accessTokenResp.AccessToken, req)
		if documentErr != nil {
			log.Error(documentErr)
			return nil, errs.FeiShuOpenApiCallError
		}

		if documentResp.Code != 0 {
			log.Error(json.ToJsonIgnoreError(documentResp))
			return nil, errs.FeiShuOpenApiCallError
		}
		total = documentResp.Data.Total
		list = append(list, documentResp.Data.DocsEntities...)
	}

	if len(list) == 0 {
		return &vo.FsDocumentListResp{
			Total: 0,
			List:  []*vo.FsDocumentData{},
		}, nil
	}

	var openIds []string
	for _, entities := range list {
		openIds = append(openIds, entities.OwnerId)
	}

	//获取用户名
	openIds = slice.SliceUniqueString(openIds)
	userMap := map[string]string{}
	if len(openIds) > 0 {
		userBatchResp, err := tenant.GetUserBatchGetV2(nil, openIds)
		if err != nil {
			log.Error(err)
			return nil, errs.FeiShuOpenApiCallError
		}
		if userBatchResp.Code != 0 {
			if userBatchResp.Code == 99991672 {
				//没有通讯录权限
				userBatchInfoResp := orgfacade.GetOrgUsersInfoByEmIds(orgvo.GetOrgUsersInfoByEmIdsReq{
					OrgId:         orgId,
					SourceChannel: sdk_const.SourceChannelFeishu,
					EmpIds:        openIds,
				})
				if userBatchInfoResp.Failure() {
					log.Error(userBatchInfoResp.Error())
					return nil, userBatchInfoResp.Error()
				}
				for _, info := range userBatchInfoResp.BaseUserInfos {
					userMap[info.OutUserId] = info.Name
				}

			} else {
				log.Error(userBatchResp.Msg)
				return nil, errs.FeiShuOpenApiCallError
			}
		} else {
			for _, user := range userBatchResp.Data.Users {
				userMap[user.OpenId] = user.Name
			}
		}
	}

	res := &vo.FsDocumentListResp{
		Total: total,
		List:  []*vo.FsDocumentData{},
	}
	for _, entities := range list {
		title := entities.Title
		if title == "" {
			title = "未命名文档"
		}
		temp := &vo.FsDocumentData{
			Title:     title,
			DocsType:  entities.DocsType,
			DocsToken: entities.DocsToken,
			OwnerID:   entities.OwnerId,
			DocsURL:   buildDocumentUrl(entities.DocsType, entities.DocsToken),
		}

		if _, ok := userMap[entities.OwnerId]; ok {
			temp.OwnerName = userMap[entities.OwnerId]
		}
		res.List = append(res.List, temp)
	}

	return res, nil
}

func buildDocumentUrl(docType string, token string) string {
	switch docType {
	case consts.FsDocumentTypeSheet:
		return fmt.Sprintf("%s/sheets/%s", consts.FsDocumentDomain, token)
	case consts.FsDocumentTypeFile:
		return fmt.Sprintf("%s/file/%s", consts.FsDocumentDomain, token)
	case consts.FsDocumentTypeDoc:
		return fmt.Sprintf("%s/docs/%s", consts.FsDocumentDomain, token)
	case consts.FsDocumentTypeDocx:
		return fmt.Sprintf("%s/docx/%s", consts.FsDocumentDomain, token)
	case consts.FsDocumentTypeBiTable:
		return fmt.Sprintf("%s/base/%s", consts.FsDocumentDomain, token)
	case consts.FsDocumentTypeMindNote:
		return fmt.Sprintf("%s/mindnotes/%s", consts.FsDocumentDomain, token)
	}
	return ""
}

// 添加飞书云文档兼容 在文件管理模块添加、附件字段添加
func AddFsResourceBatch(orgId, userId, projectId, issueId, folderId int64, data []*vo.AddIssueAttachmentFsData) ([]int64, errs.SystemErrorInfo) {
	//防止项目成员重复插入
	uid := uuid.NewUuid()
	orgIdStr := strconv.FormatInt(orgId, 10)
	relationTypeStr := strconv.Itoa(consts.FsResource)
	lockKey := consts.AddResourceLock + orgIdStr + "#" + relationTypeStr
	suc, err := cache.TryGetDistributedLock(lockKey, uid)
	if err != nil {
		log.Errorf("获取%s锁时异常 %v", lockKey, err)
		return nil, errs.TryDistributedLockError
	}
	if suc {
		defer func() {
			if _, err := cache.ReleaseDistributedLock(lockKey, uid); err != nil {
				log.Error(err)
			}
		}()
	} else {
		return nil, errs.BuildSystemErrorInfo(errs.GetDistributedLockError)
	}
	sourceType := consts.OssPolicyTypeProjectResource
	if issueId > 0 {
		sourceType = consts.OssPolicyTypeLesscodeResource
	}
	pathArr := []string{}
	originPos := []po.PpmResResource{}
	for _, datum := range data {
		host, path := str.UrlParse(datum.URL)
		pathArr = append(pathArr, path)

		suffix := util.ParseFileSuffix(datum.Title)
		//suffix := "" //飞书的后缀为空
		temp := po.PpmResResource{
			OrgId:      orgId,
			ProjectId:  projectId,
			Type:       consts.FsResource,
			Bucket:     "",
			Host:       host,
			Path:       path,
			Name:       datum.Title,
			Suffix:     suffix,
			Md5:        "",
			Size:       0,
			SourceType: sourceType,
			Creator:    userId,
		}

		//新增自动检测fileType逻辑 2019/12/30
		suffStr := strings.ToUpper(strings.TrimSpace(suffix))
		if value, ok := consts.FileTypes[suffStr]; ok {
			temp.FileType = value
		} else {
			temp.FileType = consts.FileTypeOthers
		}

		originPos = append(originPos, temp)
	}
	// 下面这段去重的逻辑去除，为了支持在文件模块不同的文件夹中展示云文档，而不至于用户以为上传失败了
	////先查看是否已经存在
	//pos := &[]po.PpmResResource{}
	//infoErr := mysql.SelectAllByCond(consts.TableResource, db.Cond{
	//	consts.TcOrgId:     orgId,
	//	consts.TcProjectId: projectId,
	//	consts.TcType:      consts.FsResource,
	//	consts.TcPath:      db.In(pathArr),
	//	consts.TcIsDelete:  consts.AppIsNoDelete,
	//}, pos)
	//if infoErr != nil {
	//	log.Error(infoErr)
	//	return nil, errs.MysqlOperateError
	//}
	//
	//result := []int64{}
	//existPath := []string{}
	//for _, resource := range *pos {
	//	existPath = append(existPath, resource.Path)
	//	result = append(result, resource.Id)
	//}
	//
	//trulyPos := []po.PpmResResource{}
	//for _, resource := range originPos {
	//	if ok, _ := slice.Contain(existPath, resource.Path); !ok {
	//		trulyPos = append(trulyPos, resource)
	//	}
	//}

	idResp, idErr := idfacade.ApplyMultiplePrimaryIdRelaxed(consts.TableResource, len(originPos))
	if idErr != nil {
		log.Error(idErr)
		return nil, idErr
	}

	var resourceIds []int64
	result := []int64{}
	resourceIdSourceTypes := make(map[int64]int)
	for i, _ := range originPos {
		originPos[i].Id = idResp.Ids[i].Id
		result = append(result, idResp.Ids[i].Id)

		resourceIds = append(resourceIds, idResp.Ids[i].Id)
		if issueId > 0 {
			resourceIdSourceTypes[idResp.Ids[i].Id] = consts.OssPolicyTypeLesscodeResource
		} else {
			resourceIdSourceTypes[idResp.Ids[i].Id] = consts.OssPolicyTypeProjectResource
		}
	}

	insertErr := mysql.BatchInsert(&po.PpmResResource{}, slice.ToSlice(originPos))
	if insertErr != nil {
		log.Error(insertErr)
		return nil, errs.MysqlOperateError
	}

	// 云文档上传至指定文件夹内，与文件夹做一个绑定关系
	idRet, errId := idfacade.ApplyMultiplePrimaryIdRelaxed(consts.TableFolderResource, len(resourceIds))
	if errId != nil {
		log.Error(errId)
		return nil, errId
	}
	folderResource := []po.PpmResFolderResource{}
	for i, resourceId := range resourceIds {
		folderResource = append(folderResource, po.PpmResFolderResource{
			Id:         idRet.Ids[i].Id,
			OrgId:      orgId,
			ResourceId: resourceId,
			FolderId:   folderId,
			Creator:    userId,
			Updator:    userId,
		})
	}
	err2 := mysql.BatchInsert(&po.PpmResFolderResource{}, slice.ToSlice(folderResource))
	if err2 != nil {
		log.Error(err2)
		return nil, errs.MysqlOperateError
	}

	// 插入resource relation 如果issueId=0，说明是从文件模块上传的
	errSys := domain.InsertResourceRelation(orgId, projectId, issueId, userId, resourceIds, resourceIdSourceTypes)
	if errSys != nil {
		log.Error(errSys)
		return nil, errSys
	}

	return result, nil
}

// 获取钉钉 最近文档列表信息
func DingDocumentList(orgId int64, userId int64, page int, size int) (*resourcevo.DingDocumentResp, errs.SystemErrorInfo) {
	baseOrgInfo, err := orgfacade.GetBaseOrgInfoRelaxed(orgId)
	if err != nil {
		log.Errorf("[DingDocumentList] orgfacade.GetBaseOrgInfoRelaxed failed, orgId:%d, userId:%d, err: %v", orgId, userId, err)
		return nil, err
	}
	client, errClient := platform_sdk.GetClient(sdk_const.SourceChannelDingTalk, baseOrgInfo.OutOrgId)
	if errClient != nil {
		return nil, errs.DingTalkClientError
	}

	// 查询出out_user_id
	userInfoBo, errSys := orgfacade.GetBaseUserInfoRelaxed(orgId, userId)
	if errSys != nil {
		log.Errorf("[DingDocumentList]orgfacade.GetBaseUserInfoRelaxed failed, orgId:%d, userId:%d, err:%v", orgId, userId, err)
		return nil, errSys
	}

	ding := client.GetOriginClient().(*dingtalk.DingTalk)
	workspaceResp, errDing := ding.GetUserWorkspaceList(userInfoBo.OutOrgUserId, true)
	if errDing != nil {
		log.Errorf("[DingDocumentList] GetUserWorkspaceList err:%v, code:%v, msg:%v", errDing, workspaceResp.Code, workspaceResp.Msg)
		return nil, errs.DingTalkOpenApiCallError
	}
	documents := make([]*resourcevo.DingDocumentData, 0, 10)
	for _, workspace := range workspaceResp.Workspaces {
		for _, doc := range workspace.RecentList {
			documents = append(documents, &resourcevo.DingDocumentData{
				DocId:   doc.NodeId,
				DocName: doc.Name,
				DocUrl:  doc.Url,
			})
		}
	}
	//docs, sdkError := client.GetDocumentList(&sdkVo.DocumentReq{
	//	Page:   page,
	//	Size:   size,
	//	UserId: userInfoBo.OutUserId,
	//})
	//if sdkError != nil {
	//	log.Errorf("[DingDocumentList]failed, orgId: %d, userId: %d, err:%v", orgId, userId, sdkError)
	//	return nil, errs.DingTalkOpenApiCallError
	//}
	//documents := make([]*resourcevo.DingDocumentData, docs.Total)
	//for _, doc := range workspaceResp.Workspaces {
	//	documents = append(documents, &resourcevo.DingDocumentData{
	//		DocId:   doc.DocId,
	//		DocName: doc.DocName,
	//		DocUrl:  doc.Url,
	//	})
	//}
	return &resourcevo.DingDocumentResp{
		Total: int64(len(documents)),
		List:  documents,
	}, nil
}

// 获取钉钉 云盘文件列表(通过云盘空间名称和类型获取)
func DingFileList(req resourcevo.DingFileListReq) (*resourcevo.DingFileListResp, errs.SystemErrorInfo) {
	baseOrgInfo, err := orgfacade.GetBaseOrgInfoRelaxed(req.OrgId)
	if err != nil {
		log.Errorf("[DingFileList] orgfacade.GetBaseOrgInfoRelaxed failed, orgId:%d, userId:%d, err: %v", req.OrgId, req.UserId, err)
		return nil, err
	}
	client, errClient := platform_sdk.GetClient(sdk_const.SourceChannelDingTalk, baseOrgInfo.OutOrgId)
	if errClient != nil {
		return nil, errs.DingTalkClientError
	}
	userInfoBo, errSys := orgfacade.GetBaseUserInfoRelaxed(req.OrgId, req.UserId)
	if errSys != nil {
		log.Errorf("[DingFileList]orgfacade.GetBaseUserInfoRelaxed failed, orgId:%d, userId:%d, err:%v", req.OrgId, req.UserId, err)
		return nil, errSys
	}
	fileList, sdkError := client.GetFileList(&sdkVo.GetFileListReq{
		Page:      req.Page,
		Size:      req.Size,
		UserId:    userInfoBo.OutUserId,
		SpaceType: req.Input.SpaceType,
		SpaceName: req.Input.SpaceName,
	})
	if sdkError != nil {
		log.Errorf("[DingFileList]failed, orgId: %d, userId: %d, err:%v", req.OrgId, req.UserId, sdkError)
		return nil, errs.DingTalkOpenApiCallError
	}

	files := make([]*resourcevo.DingFileListData, fileList.Total)

	for _, file := range fileList.List {
		files = append(files, &resourcevo.DingFileListData{
			FileType:    file.FileType,
			ContentType: file.ContentType,
			ParentId:    file.ParentId,
			FileId:      file.FileId,
			FileName:    file.FileName,
			Creator:     file.Creator,
			Modifier:    file.Modifier,
			OwnerName:   file.OwnerName,
		})
	}

	return &resourcevo.DingFileListResp{
		Total: fileList.Total,
		List:  files,
	}, nil
}

// 获取钉钉空间列表信息
func DingSpaceList(req resourcevo.DingSpaceReqVo) ([]*resourcevo.SpaceInfo, errs.SystemErrorInfo) {
	baseUserInfo := orgfacade.GetBaseUserInfo(orgvo.GetBaseUserInfoReqVo{
		OrgId:  req.OrgId,
		UserId: req.UserId,
	})
	if baseUserInfo.Failure() {
		log.Errorf("[DingSpaceList]orgfacade.GetBaseUserInfo failed:%v, orgId:%d, userId:%d",
			baseUserInfo.Error(), req.OrgId, req.UserId)
		return nil, baseUserInfo.Error()
	}
	outUserId := baseUserInfo.BaseUserInfo.OutUserId
	if outUserId == "" {
		return nil, errs.UserOutInfoNotExist
	}
	outOrgUserId := baseUserInfo.BaseUserInfo.OutOrgUserId
	spaceList, err := domain.GetDingSpaceList(outUserId, outOrgUserId, baseUserInfo.BaseUserInfo.OutOrgId, req.Input.NextToken, req.Input.MaxResults)
	if err != nil {
		log.Errorf("[DingSpaceList]domain.GetDingSpaceList err:%v, orgId:%d, userId:%d", err, req.OrgId, req.UserId)
		return nil, err
	}
	return spaceList, nil
}

// DingFileListById 只是通过spaceId和dirId获取层级列表信息
func DingFileListById(req resourcevo.DingSpaceFileReqVo) ([]*resourcevo.DingFileListData, errs.SystemErrorInfo) {
	baseUserInfo := orgfacade.GetBaseUserInfo(orgvo.GetBaseUserInfoReqVo{
		OrgId:  req.OrgId,
		UserId: req.UserId,
	})
	if baseUserInfo.Failure() {
		log.Errorf("[DingFileListById]orgfacade.GetBaseUserInfo failed:%v, orgId:%d, userId:%d",
			baseUserInfo.Error(), req.OrgId, req.UserId)
		return nil, baseUserInfo.Error()
	}
	size := req.Input.Size
	if size == 0 {
		size = 50
	}
	outOrgId := baseUserInfo.BaseUserInfo.OutOrgId
	outUserId := baseUserInfo.BaseUserInfo.OutUserId
	if outUserId == "" {
		return nil, errs.UserOutInfoNotExist
	}
	outOrgUserId := baseUserInfo.BaseUserInfo.OutOrgUserId
	files, err := domain.DingFileListById(req.OrgId, outOrgId, outUserId, outOrgUserId, req.Input.SpaceId, req.Input.DirId, size)
	if err != nil {
		log.Errorf("[DingFileListById] domain.DingFileListById err:%v", err)
		return nil, err
	}
	return files, nil
}

func GetIssueIdsByResource(orgId, projectId int64, resourceIds []int64) ([]resourcevo.GetIssueIdsByResourceIdsResp, errs.SystemErrorInfo) {
	return domain.GetIssueIdsByResource(orgId, projectId, resourceIds)
}
