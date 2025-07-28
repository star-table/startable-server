/*==============================================================*/
/* Table: ppm_bas_app_info                                      */
/*==============================================================*/
create table if not exists ppm_bas_app_info
(
   id                   bigint not null,
   name                 varchar(32) not null default '',
   code                 varchar(32) not null default '',
   secret1              varchar(32) not null default '',
   secret2              varchar(32) not null default '',
   owner                varchar(0) not null default '',
   check_status         tinyint not null default 1,
   status               tinyint not null default 1,
   creator              bigint not null default 0,
   create_time          datetime not null default CURRENT_TIMESTAMP,
   updator              bigint not null default 0,
   update_time          datetime not null default CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
   version              int not null default 1,
   is_delete            tinyint not null default 2,
   primary key (id)
);

/*==============================================================*/
/* Index: index_ppm_bas_app_info_create_time                    */
/*==============================================================*/
create index index_ppm_bas_app_info_create_time on ppm_bas_app_info
(
   create_time
);

/*==============================================================*/
/* Index: index_ppm_bas_app_info_code                           */
/*==============================================================*/
create index index_ppm_bas_app_info_code on ppm_bas_app_info
(
   code
);

/*==============================================================*/
/* Table: ppm_bas_change_log                                    */
/*==============================================================*/
create table if not exists ppm_bas_change_log
(
   id                   bigint not null,
   sys_version          varchar(64) not null default '‘’',
   feature_info         varchar(2048) not null default '',
   bug_fix_info         varchar(2048) not null default '',
   change_lnfo          varchar(2048) not null default '',
   deprecated_info      varchar(2048) not null default '',
   ext                  text,
   release_time         datetime not null default '1970-01-01 00:00:00',
   status               tinyint not null default 1,
   creator              bigint not null default 0,
   create_time          datetime not null default CURRENT_TIMESTAMP,
   updator              bigint not null default 0,
   update_time          datetime not null default CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
   version              int not null default 1,
   is_delete            tinyint not null default 2,
   primary key (id)
);

/*==============================================================*/
/* Index: index_ppm_bas_change_log_create_time                  */
/*==============================================================*/
create index index_ppm_bas_change_log_create_time on ppm_bas_change_log
(
   create_time
);

/*==============================================================*/
/* Index: index_ppm_bas_change_log_sys_version                  */
/*==============================================================*/
create index index_ppm_bas_change_log_sys_version on ppm_bas_change_log
(
   sys_version
);

/*==============================================================*/
/* Table: ppm_bas_dictionary                                    */
/*==============================================================*/
create table if not exists ppm_bas_dictionary
(
   id                   bigint not null,
   k                    varchar(64) not null default '',
   v                    varchar(4096) not null default '',
   status               tinyint not null default 1,
   creator              bigint not null default 0,
   create_time          datetime not null default CURRENT_TIMESTAMP,
   updator              bigint not null default 0,
   update_time          datetime not null default CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
   version              int not null default 1,
   is_delete            tinyint not null default 2,
   primary key (id)
);

/*==============================================================*/
/* Index: index_ppm_bas_dictionary_create_time                  */
/*==============================================================*/
create index index_ppm_bas_dictionary_create_time on ppm_bas_dictionary
(
   create_time
);

/*==============================================================*/
/* Index: index_ppm_bas_dictionary_key                          */
/*==============================================================*/
create index index_ppm_bas_dictionary_key on ppm_bas_dictionary
(
   k
);

/*==============================================================*/
/* Table: ppm_bas_object_id                                     */
/*==============================================================*/
create table if not exists ppm_bas_object_id
(
   id                   bigint not null auto_increment,
   org_id               bigint not null default 0,
   code                 varchar(64) not null default '',
   max_id               bigint not null default 0,
   step                 int not null default 100,
   creator              bigint not null default 0,
   create_time          datetime not null default CURRENT_TIMESTAMP,
   updator              bigint not null default 0,
   update_time          datetime not null default CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
   version              int not null default 1,
   is_delete            tinyint not null default 2,
   primary key (id)
);

/*==============================================================*/
/* Index: index_ppm_bas_object_id_create_time                   */
/*==============================================================*/
create index index_ppm_bas_object_id_create_time on ppm_bas_object_id
(
   create_time
);

/*==============================================================*/
/* Index: index_ppm_bas_object_id_org_id                        */
/*==============================================================*/
create index index_ppm_bas_object_id_org_id on ppm_bas_object_id
(
   org_id
);

/*==============================================================*/
/* Table: ppm_bas_pay_level                                     */
/*==============================================================*/
create table if not exists ppm_bas_pay_level
(
   id                   bigint not null,
   lang_code            varchar(64) not null default 'zh-CN',
   name                 varchar(32) not null default '',
   storage              bigint not null default 0,
   member_count         int not null default 10,
   price                bigint not null default 0,
   member_price         bigint not null default 0,
   duration             bigint not null default 0,
   is_show              tinyint not null default 1,
   sort                 int not null default 0,
   status               tinyint not null default 1,
   creator              bigint not null default 0,
   create_time          datetime not null default CURRENT_TIMESTAMP,
   updator              bigint not null default 0,
   update_time          datetime not null default CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
   version              int not null default 1,
   is_delete            tinyint not null default 2,
   primary key (id)
);

/*==============================================================*/
/* Index: index_ppm_bas_pay_level_create_time                   */
/*==============================================================*/
create index index_ppm_bas_pay_level_create_time on ppm_bas_pay_level
(
   create_time
);

/*==============================================================*/
/* Index: index_ppm_bas_pay_level_sort                          */
/*==============================================================*/
create index index_ppm_bas_pay_level_sort on ppm_bas_pay_level
(
   sort
);

/*==============================================================*/
/* Table: ppm_bas_source_channel                                */
/*==============================================================*/
create table if not exists ppm_bas_source_channel
(
   id                   bigint not null,
   code                 varchar(16) not null default '',
   name                 varchar(64) not null default '',
   remark               varchar(512) not null default '',
   status               tinyint not null default 1,
   creator              bigint not null default 0,
   create_time          datetime not null default CURRENT_TIMESTAMP,
   updator              bigint not null default 0,
   update_time          datetime not null default CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
   version              int not null default 1,
   is_delete            tinyint not null default 2,
   primary key (id)
);

/*==============================================================*/
/* Index: index_ppm_bas_source_channel_create_time              */
/*==============================================================*/
create index index_ppm_bas_source_channel_create_time on ppm_bas_source_channel
(
   create_time
);

/*==============================================================*/
/* Index: index_ppm_bas_source_channel_code                     */
/*==============================================================*/
create index index_ppm_bas_source_channel_code on ppm_bas_source_channel
(
   code
);

/*==============================================================*/
/* Table: ppm_mqs_message_queue                                 */
/*==============================================================*/
create table if not exists ppm_mqs_message_queue
(
   id                   bigint not null,
   topic                varchar(64) not null default '',
   message_key          varchar(64) not null default '',
   message              text,
   status               tinyint not null default 1,
   creator              bigint not null default 0,
   create_time          datetime not null default CURRENT_TIMESTAMP,
   updator              bigint not null default 0,
   update_time          datetime not null default CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
   version              int not null default 1,
   is_delete            tinyint not null default 2,
   primary key (id)
);

/*==============================================================*/
/* Index: index_ppm_sys_message_queue_topic_key                 */
/*==============================================================*/
create index index_ppm_sys_message_queue_topic_key on ppm_mqs_message_queue
(
   topic,
   message_key
);

/*==============================================================*/
/* Index: index_ppm_sys_message_queue_create_time               */
/*==============================================================*/
create index index_ppm_sys_message_queue_create_time on ppm_mqs_message_queue
(
   create_time
);

/*==============================================================*/
/* Table: ppm_mqs_message_queue_consumer                        */
/*==============================================================*/
create table if not exists ppm_mqs_message_queue_consumer
(
   id                   bigint not null,
   topic                varchar(64) not null default '',
   group_name           varchar(64) not null default '',
   message_id           bigint not null default 0,
   last_consumer_time   datetime not null default '1970-01-01 00:00:00',
   status               tinyint not null default 1,
   create_time          datetime not null default CURRENT_TIMESTAMP,
   update_time          datetime not null default CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
   version              int not null default 1,
   is_delete            tinyint not null default 2,
   primary key (id)
);

/*==============================================================*/
/* Index: index_ppm_mqs_message_queue_consumer_topic_group_name */
/*==============================================================*/
create index index_ppm_mqs_message_queue_consumer_topic_group_name on ppm_mqs_message_queue_consumer
(
   topic,
   group_name
);

/*==============================================================*/
/* Index: index_ppm_mqs_message_queue_consumer_create_time      */
/*==============================================================*/
create index index_ppm_mqs_message_queue_consumer_create_time on ppm_mqs_message_queue_consumer
(
   create_time
);

/*==============================================================*/
/* Table: ppm_mqs_message_queue_consumer_fail                   */
/*==============================================================*/
create table if not exists ppm_mqs_message_queue_consumer_fail
(
   id                   bigint not null,
   topic                varchar(64) not null default '',
   group_name           varchar(64) not null default '',
   message_id           bigint not null default 0,
   fail_count           int not null default 0,
   fail_time            datetime not null default '1970-01-01 00:00:00',
   status               tinyint not null default 1,
   create_time          datetime not null default CURRENT_TIMESTAMP,
   update_time          datetime not null default CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
   version              int not null default 1,
   is_delete            tinyint not null default 2,
   primary key (id)
);

/*==============================================================*/
/* Index: ppm_mqs_message_queue_consumer_fail_topic_group_name  */
/*==============================================================*/
create index ppm_mqs_message_queue_consumer_fail_topic_group_name on ppm_mqs_message_queue_consumer_fail
(
   topic,
   group_name,
   message_id
);

/*==============================================================*/
/* Index: ppm_mqs_message_queue_consumer_fail_create_time       */
/*==============================================================*/
create index ppm_mqs_message_queue_consumer_fail_create_time on ppm_mqs_message_queue_consumer_fail
(
   create_time
);

/*==============================================================*/
/* Table: ppm_orc_config                                        */
/*==============================================================*/
create table if not exists ppm_orc_config
(
   id                   bigint not null,
   org_id               bigint not null default 0,
   time_zone            varchar(32) not null default 'Asia/Shanghai',
   time_difference      varchar(8) not null default '+08:00',
   pay_level            smallint not null default 1,
   pay_start_time       datetime not null default CURRENT_TIMESTAMP,
   pay_end_time         datetime not null default '2038-01-01 00:00:00',
   web_site             varchar(256) not null default '',
   language             varchar(8) not null default 'zh-CN',
   datetime_format      varchar(32) not null default 'yyyy-MM-dd HH:mm:ss',
   password_length      tinyint not null default 6,
   password_rule        tinyint not null default 1,
   max_login_fail_count int not null default 0,
   remind_send_time     varchar(8) not null default '09:00',
   status               tinyint not null default 1,
   creator              bigint not null default 0,
   create_time          datetime not null default CURRENT_TIMESTAMP,
   updator              bigint not null default 0,
   update_time          datetime not null default CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
   version              int not null default 1,
   is_delete            tinyint not null default 2,
   primary key (id)
);

/*==============================================================*/
/* Index: index_ppm_orc_config_org_id                           */
/*==============================================================*/
create index index_ppm_orc_config_org_id on ppm_orc_config
(
   org_id
);

/*==============================================================*/
/* Index: index_ppm_orc_config_create_time                      */
/*==============================================================*/
create index index_ppm_orc_config_create_time on ppm_orc_config
(
   create_time
);

/*==============================================================*/
/* Table: ppm_orc_message_config                                */
/*==============================================================*/
create table if not exists ppm_orc_message_config
(
   id                   bigint not null,
   org_id               bigint not null default 0,
   email_status         tinyint not null default 1,
   smtp_server          varchar(128) not null default '',
   smtp_port            int not null default 25,
   smtp_user_name       varchar(128) not null default '',
   smtp_password        varchar(128) not null default '',
   email_format         varchar(8) not null default 'text',
   sender_address       varchar(128) not null default '',
   email_encode         varchar(16) not null default 'utf-8',
   status               tinyint not null default 1,
   creator              bigint not null default 0,
   create_time          datetime not null default CURRENT_TIMESTAMP,
   updator              bigint not null default 0,
   update_time          datetime not null default CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
   version              int not null default 1,
   is_delete            tinyint not null default 2,
   primary key (id)
);

/*==============================================================*/
/* Index: index_ppm_orc_message_config_create_time              */
/*==============================================================*/
create index index_ppm_orc_message_config_create_time on ppm_orc_message_config
(
   create_time
);

/*==============================================================*/
/* Index: index_ppm_orc_message_config_org_id                   */
/*==============================================================*/
create index index_ppm_orc_message_config_org_id on ppm_orc_message_config
(
   org_id
);

/*==============================================================*/
/* Table: ppm_org_department                                    */
/*==============================================================*/
create table if not exists ppm_org_department
(
   id                   bigint not null,
   org_id               bigint not null default 0,
   name                 varchar(64) not null default '',
   code                 varchar(64) not null default '',
   parent_id            bigint not null default 0,
   sort                 int not null default 0,
   is_hide              tinyint not null default 2,
   source_channel       varchar(16) not null default '',
   status               tinyint not null default 1,
   creator              bigint not null default 0,
   create_time          datetime not null default CURRENT_TIMESTAMP,
   updator              bigint not null default 0,
   update_time          datetime not null default CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
   version              int not null default 1,
   is_delete            tinyint not null default 2,
   primary key (id)
);

/*==============================================================*/
/* Index: index_ppm_org_department_org_id                       */
/*==============================================================*/
create index index_ppm_org_department_org_id on ppm_org_department
(
   org_id
);

/*==============================================================*/
/* Index: index_ppm_org_department_code                         */
/*==============================================================*/
create index index_ppm_org_department_code on ppm_org_department
(
   code
);

/*==============================================================*/
/* Index: index_ppm_org_department_parent_id                    */
/*==============================================================*/
create index index_ppm_org_department_parent_id on ppm_org_department
(
   parent_id
);

/*==============================================================*/
/* Index: index_ppm_org_department_create_time                  */
/*==============================================================*/
create index index_ppm_org_department_create_time on ppm_org_department
(
   create_time
);

/*==============================================================*/
/* Table: ppm_org_department_out_info                           */
/*==============================================================*/
create table if not exists ppm_org_department_out_info
(
   id                   bigint not null,
   org_id               bigint not null default 0,
   department_id        bigint not null default 0,
   source_channel       varchar(16) not null default '',
   out_org_department_id varchar(64) not null default '',
   out_org_department_code varchar(64) not null default '',
   name                 varchar(64) not null default '',
   out_org_department_parent_id varchar(64) not null default '',
   status               tinyint not null default 1,
   creator              bigint not null default 0,
   create_time          datetime not null default CURRENT_TIMESTAMP,
   updator              bigint not null default 0,
   update_time          datetime not null default CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
   version              int not null default 1,
   is_delete            tinyint not null default 2,
   primary key (id)
);

/*==============================================================*/
/* Index: index_ppm_org_department_out_info_org_id              */
/*==============================================================*/
create index index_ppm_org_department_out_info_org_id on ppm_org_department_out_info
(
   org_id
);

/*==============================================================*/
/* Index: index_ppm_org_department_out_info_department_id       */
/*==============================================================*/
create index index_ppm_org_department_out_info_department_id on ppm_org_department_out_info
(
   department_id
);

/*================================================================*/
/* Index: index_ppm_org_department_out_info_out_org_department_id */
/*================================================================*/
create index index_ppm_org_department_out_info_out_org_department_id on ppm_org_department_out_info
(
   out_org_department_id
);

/*==================================================================*/
/* Index: index_ppm_org_department_out_info_out_org_department_code */
/*==================================================================*/
create index index_ppm_org_department_out_info_out_org_department_code on ppm_org_department_out_info
(
   out_org_department_code
);

/*==============================================================*/
/* Index: index_ppm_org_department_out_info_create_time         */
/*==============================================================*/
create index index_ppm_org_department_out_info_create_time on ppm_org_department_out_info
(
   create_time
);

/*=======================================================================*/
/* Index: index_ppm_org_department_out_info_out_org_department_parent_id */
/*=======================================================================*/
create index index_ppm_org_department_out_info_out_org_department_parent_id on ppm_org_department_out_info
(
   out_org_department_parent_id
);

/*==============================================================*/
/* Table: ppm_org_organization                                  */
/*==============================================================*/
create table if not exists ppm_org_organization
(
   id                   bigint not null,
   name                 varchar(256) not null default '',
   web_site             varchar(512) not null default '',
   industry_id          bigint not null default 0,
   scale                varchar(32) not null default '',
   source_channel       varchar(16) not null default '',
   country_id           bigint not null default 0,
   province_id          bigint not null default 0,
   city_id              bigint not null default 0,
   address              varchar(256) not null default '',
   logo_url             varchar(512) not null default '',
   resource_id          bigint not null default 0,
   owner                bigint not null default 0,
   is_authenticated     tinyint not null default 1,
   remark               varchar(512) not null default '',
   init_status          tinyint not null default 1,
   init_version         int not null default 1,
   status               tinyint not null default 1,
   is_show              tinyint not null default 1,
   creator              bigint not null default 0,
   create_time          datetime not null default CURRENT_TIMESTAMP,
   updator              bigint not null default 0,
   update_time          datetime not null default CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
   version              int not null default 1,
   is_delete            tinyint not null default 2,
   primary key (id)
);

/*==============================================================*/
/* Index: index_ppm_org_organization_create_time                */
/*==============================================================*/
create index index_ppm_org_organization_create_time on ppm_org_organization
(
   create_time
);

/*==============================================================*/
/* Index: index_ppm_org_organization_city_id                    */
/*==============================================================*/
create index index_ppm_org_organization_city_id on ppm_org_organization
(
   city_id
);

/*==============================================================*/
/* Index: index_ppm_org_organization_province_id                */
/*==============================================================*/
create index index_ppm_org_organization_province_id on ppm_org_organization
(
   province_id
);

/*==============================================================*/
/* Index: index_ppm_org_organization_country_id                 */
/*==============================================================*/
create index index_ppm_org_organization_country_id on ppm_org_organization
(
   country_id
);

/*==============================================================*/
/* Table: ppm_org_organization_out_info                         */
/*==============================================================*/
create table if not exists ppm_org_organization_out_info
(
   id                   bigint not null,
   org_id               bigint not null default 0,
   out_org_id           varchar(64) not null default '',
   source_channel       varchar(16) not null default '',
   name                 varchar(64) not null default '',
   industry             varchar(64) not null default '',
   is_authenticated     tinyint not null default 1,
   auth_ticket          varchar(256) not null default '',
   auth_level           varchar(32) not null default '',
   status               tinyint not null default 1,
   creator              bigint not null default 0,
   create_time          datetime not null default CURRENT_TIMESTAMP,
   updator              bigint not null default 0,
   update_time          datetime not null default CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
   version              int not null default 1,
   is_delete            tinyint not null default 2,
   primary key (id)
);

/*==============================================================*/
/* Index: index_ppm_org_organization_out_info_org_id            */
/*==============================================================*/
create index index_ppm_org_organization_out_info_org_id on ppm_org_organization_out_info
(
   org_id
);

/*==============================================================*/
/* Index: index_ppm_org_organization_out_info_out_org_id        */
/*==============================================================*/
create index index_ppm_org_organization_out_info_out_org_id on ppm_org_organization_out_info
(
   out_org_id
);

/*==============================================================*/
/* Index: index_ppm_org_organization_out_info_create_time       */
/*==============================================================*/
create index index_ppm_org_organization_out_info_create_time on ppm_org_organization_out_info
(
   create_time
);

/*==============================================================*/
/* Table: ppm_org_user                                          */
/*==============================================================*/
create table if not exists ppm_org_user
(
   id                   bigint not null,
   org_id               bigint not null default 0,
   name                 varchar(64) not null default '',
   login_name           varchar(32) not null,
   login_name_edit_count int not null default 0,
   email                varchar(128) not null default '',
   mobile               varchar(16) not null default '',
   avatar               varchar(1024) not null default '',
   birthday             datetime not null default '1970-01-01 00:00:00',
   sex                  tinyint not null default 99,
   password             varchar(32) not null default '',
   password_salt        varchar(32) not null default '',
   source_channel       varchar(16) not null default '',
   language             varchar(8) not null default 'zh-CN',
   motto                varchar(512) not null default '',
   last_login_ip        varchar(64) not null default '',
   last_login_time      datetime not null default '1970-01-01 00:00:00',
   login_fail_count     int not null default 0,
   status               tinyint not null default 1,
   creator              bigint not null default 0,
   create_time          datetime not null default CURRENT_TIMESTAMP,
   updator              bigint not null default 0,
   update_time          datetime not null default CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
   version              int not null default 1,
   is_delete            tinyint not null default 2,
   primary key (id)
);

/*==============================================================*/
/* Index: index_ppm_org_user_org_id                             */
/*==============================================================*/
create index index_ppm_org_user_org_id on ppm_org_user
(
   org_id
);

/*==============================================================*/
/* Index: index_ppm_org_user_mobile                             */
/*==============================================================*/
create index index_ppm_org_user_mobile on ppm_org_user
(
   mobile
);

/*==============================================================*/
/* Index: index_ppm_org_user_login_name                         */
/*==============================================================*/
create index index_ppm_org_user_login_name on ppm_org_user
(
   login_name
);

/*==============================================================*/
/* Index: index_ppm_org_user_email                              */
/*==============================================================*/
create index index_ppm_org_user_email on ppm_org_user
(
   email
);

/*==============================================================*/
/* Index: index_ppm_org_user_create_time                        */
/*==============================================================*/
create index index_ppm_org_user_create_time on ppm_org_user
(
   create_time
);

/*==============================================================*/
/* Table: ppm_org_user_config                                   */
/*==============================================================*/
create table if not exists ppm_org_user_config
(
   id                   bigint not null,
   org_id               bigint not null default 0,
   user_id              bigint not null default 0,
   daily_report_message_status tinyint not null default 2,
   owner_range_status   tinyint not null default 1,
   participant_range_status tinyint not null default 1,
   attention_range_status tinyint not null default 1,
   create_range_status  tinyint not null default 2,
   remind_message_status tinyint not null default 2,
   comment_at_message_status tinyint not null default 1,
   modify_message_status tinyint not null default 1,
   relation_message_status tinyint not null default 2,
   daily_project_report_message_status tinyint not null default 2,
   default_project_id bigint not null default 0,
   ext                  varchar(4096) not null default '',
   creator              bigint not null default 0,
   create_time          datetime not null default CURRENT_TIMESTAMP,
   updator              bigint not null default 0,
   update_time          datetime not null default CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
   version              int not null default 1,
   is_delete            tinyint not null default 2,
   primary key (id)
);

/*==============================================================*/
/* Index: index_ppm_org_user_config_user_id                     */
/*==============================================================*/
create index index_ppm_org_user_config_user_id on ppm_org_user_config
(
   user_id
);

/*==============================================================*/
/* Index: index_ppm_org_user_config_org_id                      */
/*==============================================================*/
create index index_ppm_org_user_config_org_id on ppm_org_user_config
(
   org_id
);

/*==============================================================*/
/* Index: index_ppm_org_user_config_create_time                 */
/*==============================================================*/
create index index_ppm_org_user_config_create_time on ppm_org_user_config
(
   create_time
);

/*==============================================================*/
/* Table: ppm_org_user_department                               */
/*==============================================================*/
create table if not exists ppm_org_user_department
(
   id                   bigint not null,
   org_id               bigint not null default 0,
   user_id              bigint not null default 0,
   department_id        bigint not null default 0,
   is_leader            tinyint not null default 2,
   creator              bigint not null default 0,
   create_time          datetime not null default CURRENT_TIMESTAMP,
   updator              bigint not null default 0,
   update_time          datetime not null default CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
   version              int not null default 1,
   is_delete            tinyint not null default 2,
   primary key (id)
);

/*==============================================================*/
/* Index: index_ppm_org_user_department_org_id                  */
/*==============================================================*/
create index index_ppm_org_user_department_org_id on ppm_org_user_department
(
   org_id
);

/*==============================================================*/
/* Index: index_ppm_org_user_department_create_time             */
/*==============================================================*/
create index index_ppm_org_user_department_create_time on ppm_org_user_department
(
   create_time
);

/*==============================================================*/
/* Index: index_ppm_org_user_department_user_id                 */
/*==============================================================*/
create index index_ppm_org_user_department_user_id on ppm_org_user_department
(
   user_id
);

/*==============================================================*/
/* Index: index_ppm_org_user_department_department_id           */
/*==============================================================*/
create index index_ppm_org_user_department_department_id on ppm_org_user_department
(
   department_id
);

/*==============================================================*/
/* Table: ppm_org_user_organization                             */
/*==============================================================*/
create table if not exists ppm_org_user_organization
(
   id                   bigint not null,
   org_id               bigint not null default 0,
   user_id              bigint not null default 0,
   check_status         tinyint not null default 1,
   use_status           tinyint not null default 2,
   status               tinyint not null default 1,
   creator              bigint not null default 0,
   create_time          datetime not null default CURRENT_TIMESTAMP,
   updator              bigint not null default 0,
   update_time          datetime not null default CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
   version              int not null default 1,
   is_delete            tinyint not null default 2,
   primary key (id)
);

/*==============================================================*/
/* Index: index_ppm_org_user_organization_org_id                */
/*==============================================================*/
create index index_ppm_org_user_organization_org_id on ppm_org_user_organization
(
   org_id
);

/*==============================================================*/
/* Index: index_ppm_org_user_organization_create_time           */
/*==============================================================*/
create index index_ppm_org_user_organization_create_time on ppm_org_user_organization
(
   create_time
);

/*==============================================================*/
/* Index: index_ppm_org_user_organization_user_id               */
/*==============================================================*/
create index index_ppm_org_user_organization_user_id on ppm_org_user_organization
(
   user_id
);

/*==============================================================*/
/* Table: ppm_org_user_out_info                                 */
/*==============================================================*/
create table if not exists ppm_org_user_out_info
(
   id                   bigint not null,
   org_id               bigint not null default 0,
   user_id              bigint not null default 0,
   source_channel       varchar(16) not null default '',
   out_org_user_id      varchar(64) not null default '',
   out_user_id          varchar(64) not null default '',
   name                 varchar(64) not null default '',
   avatar               varchar(1024) not null default '',
   is_active            tinyint not null default 1,
   job_number           varchar(32) not null default '',
   status               tinyint not null default 1,
   creator              bigint not null default 0,
   create_time          datetime not null default CURRENT_TIMESTAMP,
   updator              bigint not null default 0,
   update_time          datetime not null default CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
   version              int not null default 1,
   is_delete            tinyint not null default 2,
   primary key (id)
);

/*==============================================================*/
/* Index: index_ppm_org_user_out_info_org_id                    */
/*==============================================================*/
create index index_ppm_org_user_out_info_org_id on ppm_org_user_out_info
(
   org_id
);

/*==============================================================*/
/* Index: index_ppm_org_user_out_info_user_id                   */
/*==============================================================*/
create index index_ppm_org_user_out_info_user_id on ppm_org_user_out_info
(
   user_id
);

/*==============================================================*/
/* Index: index_ppm_org_user_out_info_out_org_user_id           */
/*==============================================================*/
create index index_ppm_org_user_out_info_out_org_user_id on ppm_org_user_out_info
(
   out_org_user_id
);

/*==============================================================*/
/* Index: index_ppm_org_user_out_info_out_user_id               */
/*==============================================================*/
create index index_ppm_org_user_out_info_out_user_id on ppm_org_user_out_info
(
   out_user_id
);

/*==============================================================*/
/* Index: index_ppm_org_user_out_info_create_time               */
/*==============================================================*/
create index index_ppm_org_user_out_info_create_time on ppm_org_user_out_info
(
   create_time
);

/*==============================================================*/
/* Table: ppm_pri_issue                                         */
/*==============================================================*/
create table if not exists ppm_pri_issue
(
   id                   bigint not null,
   org_id               bigint not null default 0,
   code                 varchar(64) not null default '',
   project_id           bigint not null default 0,
   project_object_type_id bigint not null default 0,
   title                varchar(512) not null default '',
   owner                bigint not null default 0,
   priority_id          bigint not null default 0,
   source_id            bigint not null default 0,
   issus_object_type_id bigint not null default 0,
   plan_start_time      datetime not null default '1970-01-01 00:00:00',
   plan_end_time        datetime not null default '1970-01-01 00:00:00',
   start_time           datetime not null default '1970-01-01 00:00:00',
   end_time             datetime not null default '1970-01-01 00:00:00',
   plan_work_hour       int not null default -1,
   iteration_id         bigint not null default 0,
   version_id           bigint not null default 0,
   module_id            bigint not null default 0,
   parent_id            bigint not null default 0,
   status               bigint not null default 0,
   creator              bigint not null default 0,
   create_time          datetime not null default CURRENT_TIMESTAMP,
   updator              bigint not null default 0,
   update_time          datetime not null default CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
   version              int not null default 1,
   is_delete            tinyint not null default 2,
   primary key (id)
);

/*==============================================================*/
/* Index: index_ppm_pri_issue_org_id                            */
/*==============================================================*/
create index index_ppm_pri_issue_org_id on ppm_pri_issue
(
   org_id
);

/*==============================================================*/
/* Index: index_ppm_pri_issue_project_id_project_objct_type_id  */
/*==============================================================*/
create index index_ppm_pri_issue_project_id_project_objct_type_id on ppm_pri_issue
(
   project_id,
   project_object_type_id
);

/*==============================================================*/
/* Index: index_ppm_pri_issue_code                              */
/*==============================================================*/
create index index_ppm_pri_issue_code on ppm_pri_issue
(
   code
);

/*==============================================================*/
/* Index: index_ppm_pri_issue_handler                           */
/*==============================================================*/
create index index_ppm_pri_issue_handler on ppm_pri_issue
(
   owner
);

/*==============================================================*/
/* Index: index_ppm_pri_issue_parent_id                         */
/*==============================================================*/
create index index_ppm_pri_issue_parent_id on ppm_pri_issue
(
   parent_id
);

/*==============================================================*/
/* Index: index_ppm_pri_issue_status                            */
/*==============================================================*/
create index index_ppm_pri_issue_status on ppm_pri_issue
(
   status
);

/*==============================================================*/
/* Index: index_ppm_pri_issue_iteration_id                      */
/*==============================================================*/
create index index_ppm_pri_issue_iteration_id on ppm_pri_issue
(
   iteration_id
);

/*==============================================================*/
/* Index: index_ppm_pri_issue_creator                           */
/*==============================================================*/
create index index_ppm_pri_issue_creator on ppm_pri_issue
(
   creator
);

/*==============================================================*/
/* Index: index_ppm_pri_issue_create_time                       */
/*==============================================================*/
create index index_ppm_pri_issue_create_time on ppm_pri_issue
(
   create_time
);

/*==============================================================*/
/* Index: index_ppm_pri_issue_end_time                          */
/*==============================================================*/
create index index_ppm_pri_issue_end_time on ppm_pri_issue
(
   end_time
);

/*==============================================================*/
/* Index: index_ppm_pri_issue_start_time                        */
/*==============================================================*/
create index index_ppm_pri_issue_start_time on ppm_pri_issue
(
   start_time
);

/*==============================================================*/
/* Index: index_ppm_pri_issue_plan_start_time                   */
/*==============================================================*/
create index index_ppm_pri_issue_plan_start_time on ppm_pri_issue
(
   plan_start_time
);

/*==============================================================*/
/* Index: index_ppm_pri_issue_plan_end_time                     */
/*==============================================================*/
create index index_ppm_pri_issue_plan_end_time on ppm_pri_issue
(
   plan_end_time
);

/*==============================================================*/
/* Table: ppm_pri_issue_detail                                  */
/*==============================================================*/
create table if not exists ppm_pri_issue_detail
(
   id                   bigint not null,
   org_id               bigint not null default 0,
   issue_id             bigint not null default 0,
   project_id           bigint not null default 0,
   story_point          int not null default -1,
   tags                 varchar(1024) not null default '',
   remark               text,
   status               bigint not null default 0,
   creator              bigint not null default 0,
   create_time          datetime not null default CURRENT_TIMESTAMP,
   updator              bigint not null default 0,
   update_time          datetime not null default CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
   version              int not null default 1,
   is_delete            tinyint not null default 2,
   primary key (id)
);

/*==============================================================*/
/* Index: index_ppm_pri_issue_detail_org_id                     */
/*==============================================================*/
create index index_ppm_pri_issue_detail_org_id on ppm_pri_issue_detail
(
   org_id
);

/*==============================================================*/
/* Index: index_ppm_pri_issue_detail_issue_id                   */
/*==============================================================*/
create index index_ppm_pri_issue_detail_issue_id on ppm_pri_issue_detail
(
   issue_id
);

/*==============================================================*/
/* Index: index_ppm_pri_issue_detail_project_id                 */
/*==============================================================*/
create index index_ppm_pri_issue_detail_project_id on ppm_pri_issue_detail
(
   project_id
);

/*==============================================================*/
/* Index: index_ppm_pri_issue_detail_create_time                */
/*==============================================================*/
create index index_ppm_pri_issue_detail_create_time on ppm_pri_issue_detail
(
   create_time
);

/*==============================================================*/
/* Table: ppm_pri_issue_relation                                */
/*==============================================================*/
create table if not exists ppm_pri_issue_relation
(
   id                   bigint not null,
   org_id               bigint not null default 0,
   project_id           bigint not null default 0,
   issue_id             bigint not null default 0,
   relation_id          bigint not null default 0,
   relation_type        tinyint not null default 1,
   creator              bigint not null default 0,
   create_time          datetime not null default CURRENT_TIMESTAMP,
   updator              bigint not null default 0,
   update_time          datetime not null default CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
   version              int not null default 1,
   is_delete            tinyint not null default 2
);

/*==============================================================*/
/* Index: index_ppm_pri_issue_relation_project_id               */
/*==============================================================*/
create index index_ppm_pri_issue_relation_project_id on ppm_pri_issue_relation
(
   project_id
);

/*==============================================================*/
/* Index: index_ppm_pri_issue_relation_org_id                   */
/*==============================================================*/
create index index_ppm_pri_issue_relation_org_id on ppm_pri_issue_relation
(
   org_id
);

/*==============================================================*/
/* Index: index_ppm_pri_issue_relation_issue_id                 */
/*==============================================================*/
create index index_ppm_pri_issue_relation_issue_id on ppm_pri_issue_relation
(
   issue_id
);

/*==============================================================*/
/* Index: index_ppm_pri_issue_relation_relation_id              */
/*==============================================================*/
create index index_ppm_pri_issue_relation_relation_id on ppm_pri_issue_relation
(
   relation_id
);

/*==============================================================*/
/* Index: index_ppm_pri_issue_relation_create_time              */
/*==============================================================*/
create index index_ppm_pri_issue_relation_create_time on ppm_pri_issue_relation
(
   create_time
);

/*==============================================================*/
/* Table: ppm_pri_iteration                                     */
/*==============================================================*/
create table if not exists ppm_pri_iteration
(
   id                   bigint not null,
   org_id               bigint not null default 0,
   project_id           bigint not null default 0,
   name                 varchar(64) not null default '',
   owner                bigint not null default 0,
   version_id           bigint not null default 0,
   plan_start_time      datetime not null default '1970-01-01 00:00:00',
   plan_end_time        datetime not null default '1970-01-01 00:00:00',
   plan_work_hour       int not null default -1,
   story_point          int not null default -1,
   remark               text,
   status               bigint not null default 0,
   creator              bigint not null default 0,
   create_time          datetime not null default CURRENT_TIMESTAMP,
   updator              bigint not null default 0,
   update_time          datetime not null default CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
   version              int not null default 1,
   is_delete            tinyint not null default 2,
   primary key (id)
);

/*==============================================================*/
/* Index: index_ppm_pri_iteration_org_id                        */
/*==============================================================*/
create index index_ppm_pri_iteration_org_id on ppm_pri_iteration
(
   org_id
);

/*==============================================================*/
/* Index: index_ppm_pri_iteration_project_id                    */
/*==============================================================*/
create index index_ppm_pri_iteration_project_id on ppm_pri_iteration
(
   project_id
);

/*==============================================================*/
/* Index: index_ppm_pri_iteration_owner                         */
/*==============================================================*/
create index index_ppm_pri_iteration_owner on ppm_pri_iteration
(
   owner
);

/*==============================================================*/
/* Index: index_ppm_pri_iteration_version_id                    */
/*==============================================================*/
create index index_ppm_pri_iteration_version_id on ppm_pri_iteration
(
   version_id
);

/*==============================================================*/
/* Index: index_ppm_pri_iteration_creator                       */
/*==============================================================*/
create index index_ppm_pri_iteration_creator on ppm_pri_iteration
(
   creator
);

/*==============================================================*/
/* Index: index_ppm_pri_iteration_create_time                   */
/*==============================================================*/
create index index_ppm_pri_iteration_create_time on ppm_pri_iteration
(
   create_time
);

/*==============================================================*/
/* Table: ppm_pro_application                                   */
/*==============================================================*/
create table if not exists ppm_pro_application
(
   id                   bigint not null,
   org_id               bigint not null default 0,
   project_id           bigint not null default 0,
   name                 varchar(128) not null default '',
   package              varchar(256) not null default '',
   icon                 varchar(256) not null default '',
   owner                bigint not null default 0,
   remark               varchar(512) not null default '',
   status               tinyint not null default 1,
   creator              bigint not null default 0,
   create_time          datetime not null default CURRENT_TIMESTAMP,
   updator              bigint not null default 0,
   update_time          datetime not null default CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
   version              int not null default 1,
   is_delete            tinyint not null default 2,
   primary key (id)
);

/*==============================================================*/
/* Index: index_ppm_pro_application_org_id                      */
/*==============================================================*/
create index index_ppm_pro_application_org_id on ppm_pro_application
(
   org_id
);

/*==============================================================*/
/* Index: index_ppm_pro_application_project_id                  */
/*==============================================================*/
create index index_ppm_pro_application_project_id on ppm_pro_application
(
   project_id
);

/*==============================================================*/
/* Index: index_ppm_pro_project_member_create_time              */
/*==============================================================*/
create index index_ppm_pro_project_member_create_time on ppm_pro_application
(
   create_time
);

/*==============================================================*/
/* Table: ppm_pro_application_version                           */
/*==============================================================*/
create table if not exists ppm_pro_application_version
(
   id                   bigint not null,
   org_id               bigint not null default 0,
   app_id               bigint not null default 0,
   app_version          varchar(256) not null default '',
   owner                bigint not null default 0,
   status               tinyint not null default 1,
   creator              bigint not null default 0,
   create_time          datetime not null default CURRENT_TIMESTAMP,
   updator              bigint not null default 0,
   update_time          datetime not null default CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
   version              int not null default 1,
   is_delete            tinyint not null default 2,
   primary key (id)
);

/*==============================================================*/
/* Index: index_ppm_pro_application_version_org_id              */
/*==============================================================*/
create index index_ppm_pro_application_version_org_id on ppm_pro_application_version
(
   org_id
);

/*==============================================================*/
/* Index: index_ppm_pro_application_version_app_id              */
/*==============================================================*/
create index index_ppm_pro_application_version_app_id on ppm_pro_application_version
(
   app_id
);

/*==============================================================*/
/* Index: index_ppm_pro_project_version_create_time             */
/*==============================================================*/
create index index_ppm_pro_project_version_create_time on ppm_pro_application_version
(
   create_time
);

/*==============================================================*/
/* Table: ppm_pro_project                                       */
/*==============================================================*/
create table if not exists ppm_pro_project
(
   id                   bigint not null,
   org_id               bigint not null default 0,
   code                 varchar(64) not null default '',
   name                 varchar(256) not null default '',
   pre_code             varchar(16) not null default '',
   owner                bigint not null default 0,
   project_type_id      bigint not null default 1,
   priority_id          bigint not null default 0,
   plan_start_time      datetime not null default '1970-01-01 00:00:00',
   plan_end_time        datetime not null default '1970-01-01 00:00:00',
   public_status        tinyint not null default 1,
   resource_id          bigint not null default 0,
   is_filing            tinyint not null default 2,
   remark               varchar(512) not null default '',
   status               bigint not null default 0,
   creator              bigint not null default 0,
   create_time          datetime not null default CURRENT_TIMESTAMP,
   updator              bigint not null default 0,
   update_time          datetime not null default CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
   version              int not null default 1,
   is_delete            tinyint not null default 2,
   primary key (id)
);

/*==============================================================*/
/* Index: index_ppm_pro_project_org_id                          */
/*==============================================================*/
create index index_ppm_pro_project_org_id on ppm_pro_project
(
   org_id
);

/*==============================================================*/
/* Index: index_ppm_pro_project_pre_code                        */
/*==============================================================*/
create index index_ppm_pro_project_pre_code on ppm_pro_project
(
   pre_code
);

/*==============================================================*/
/* Index: index_ppm_pro_project_code                            */
/*==============================================================*/
create index index_ppm_pro_project_code on ppm_pro_project
(
   code
);

/*==============================================================*/
/* Index: index_ppm_pro_project_owner                           */
/*==============================================================*/
create index index_ppm_pro_project_owner on ppm_pro_project
(
   owner
);

/*==============================================================*/
/* Index: index_ppm_pro_project_name                            */
/*==============================================================*/
create index index_ppm_pro_project_name on ppm_pro_project
(
   name
);

/*==============================================================*/
/* Index: index_ppm_pro_project_create_time                     */
/*==============================================================*/
create index index_ppm_pro_project_create_time on ppm_pro_project
(
   create_time
);

/*==============================================================*/
/* Index: index_ppm_pro_project_plan_start_time                 */
/*==============================================================*/
create index index_ppm_pro_project_plan_start_time on ppm_pro_project
(
   plan_start_time
);

/*==============================================================*/
/* Index: index_ppm_pro_project_plan_end_time                   */
/*==============================================================*/
create index index_ppm_pro_project_plan_end_time on ppm_pro_project
(
   plan_end_time
);

/*==============================================================*/
/* Table: ppm_pro_project_detail                                */
/*==============================================================*/
create table if not exists ppm_pro_project_detail
(
   id                   bigint not null,
   org_id               bigint not null default 0,
   project_id           bigint not null default 0,
   notice               varchar(4096) not null default '',
   is_enable_work_hours tinyint not null default 1,
   creator              bigint not null default 0,
   create_time          datetime not null default CURRENT_TIMESTAMP,
   updator              bigint not null default 0,
   update_time          datetime not null default CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
   version              int not null default 1,
   is_delete            tinyint not null default 2,
   primary key (id)
);

/*==============================================================*/
/* Index: index_ppm_pro_project_detail_org_id                   */
/*==============================================================*/
create index index_ppm_pro_project_detail_org_id on ppm_pro_project_detail
(
   org_id
);

/*==============================================================*/
/* Index: index_ppm_pro_project_detail_project_id               */
/*==============================================================*/
create index index_ppm_pro_project_detail_project_id on ppm_pro_project_detail
(
   project_id
);

/*==============================================================*/
/* Index: index_ppm_pro_project_detail_create_time              */
/*==============================================================*/
create index index_ppm_pro_project_detail_create_time on ppm_pro_project_detail
(
   create_time
);

/*==============================================================*/
/* Table: ppm_pro_project_module                                */
/*==============================================================*/
create table if not exists ppm_pro_project_module
(
   id                   bigint not null,
   org_id               bigint not null default 0,
   project_id           bigint not null default 0,
   name                 varchar(64) not null default '',
   code                 varchar(32) not null default '',
   owner                bigint not null default 0,
   remark               varchar(512) not null default '',
   status               tinyint not null default 1,
   creator              bigint not null default 0,
   create_time          datetime not null default CURRENT_TIMESTAMP,
   updator              bigint not null default 0,
   update_time          datetime not null default CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
   version              int not null default 1,
   is_delete            tinyint not null default 2,
   primary key (id)
);

/*==============================================================*/
/* Index: index_ppm_pro_project_module_org_id                   */
/*==============================================================*/
create index index_ppm_pro_project_module_org_id on ppm_pro_project_module
(
   org_id
);

/*==============================================================*/
/* Index: index_ppm_pro_project_module_project_id               */
/*==============================================================*/
create index index_ppm_pro_project_module_project_id on ppm_pro_project_module
(
   project_id
);

/*==============================================================*/
/* Index: index_ppm_pro_project_module_owner                    */
/*==============================================================*/
create index index_ppm_pro_project_module_owner on ppm_pro_project_module
(
   owner
);

/*==============================================================*/
/* Index: index_ppm_pro_project_module_create_time              */
/*==============================================================*/
create index index_ppm_pro_project_module_create_time on ppm_pro_project_module
(
   create_time
);

/*==============================================================*/
/* Table: ppm_pro_project_relation                              */
/*==============================================================*/
create table if not exists ppm_pro_project_relation
(
   id                   bigint not null,
   org_id               bigint not null default 0,
   project_id           bigint not null default 0,
   team_id              bigint not null default 0,
   relation_id          bigint not null default 0,
   relation_type        tinyint not null default 1,
   status               tinyint not null default 1,
   creator              bigint not null default 0,
   create_time          datetime not null default CURRENT_TIMESTAMP,
   updator              bigint not null default 0,
   update_time          datetime not null default CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
   version              int not null default 1,
   is_delete            tinyint not null default 2,
   primary key (id)
);

/*==============================================================*/
/* Index: index_ppm_pro_project_relation_org_id                 */
/*==============================================================*/
create index index_ppm_pro_project_relation_org_id on ppm_pro_project_relation
(
   org_id
);

/*==============================================================*/
/* Index: index_ppm_pro_project_relation_project_id             */
/*==============================================================*/
create index index_ppm_pro_project_relation_project_id on ppm_pro_project_relation
(
   project_id
);

/*==============================================================*/
/* Index: index_ppm_pro_project_relation_relation_id            */
/*==============================================================*/
create index index_ppm_pro_project_relation_relation_id on ppm_pro_project_relation
(
   relation_id
);

/*==============================================================*/
/* Index: index_ppm_pro_project_relation_create_time            */
/*==============================================================*/
create index index_ppm_pro_project_relation_create_time on ppm_pro_project_relation
(
   create_time
);

/*==============================================================*/
/* Table: ppm_pro_project_version                               */
/*==============================================================*/
create table if not exists ppm_pro_project_version
(
   id                   bigint not null,
   org_id               bigint not null default 0,
   project_id           bigint not null default 0,
   name                 varchar(64) not null default '',
   code                 varchar(32) not null default '',
   owner                bigint not null default 0,
   remark               varchar(512) not null default '',
   status               tinyint not null default 1,
   creator              bigint not null default 0,
   create_time          datetime not null default CURRENT_TIMESTAMP,
   updator              bigint not null default 0,
   update_time          datetime not null default CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
   version              int not null default 1,
   is_delete            tinyint not null default 2,
   primary key (id)
);

/*==============================================================*/
/* Index: index_ppm_pro_project_version_org_id                  */
/*==============================================================*/
create index index_ppm_pro_project_version_org_id on ppm_pro_project_version
(
   org_id
);

/*==============================================================*/
/* Index: index_ppm_pro_project_version_project_id              */
/*==============================================================*/
create index index_ppm_pro_project_version_project_id on ppm_pro_project_version
(
   project_id
);

/*==============================================================*/
/* Index: index_ppm_pro_project_version_owner                   */
/*==============================================================*/
create index index_ppm_pro_project_version_owner on ppm_pro_project_version
(
   owner
);

/*==============================================================*/
/* Index: index_ppm_pro_project_version_create_time             */
/*==============================================================*/
create index index_ppm_pro_project_version_create_time on ppm_pro_project_version
(
   create_time
);

/*==============================================================*/
/* Table: ppm_prs_issue_object_type                             */
/*==============================================================*/
create table if not exists ppm_prs_issue_object_type
(
   id                   bigint not null,
   org_id               bigint not null default 0,
   lang_code            varchar(64) not null default '',
   name                 varchar(64) not null default '',
   sort                 int not null default 0,
   project_object_type_id bigint not null default 0,
   remark               varchar(512) not null default '',
   status               tinyint not null default 1,
   creator              bigint not null default 0,
   create_time          datetime not null default CURRENT_TIMESTAMP,
   updator              bigint not null default 0,
   update_time          datetime not null default CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
   version              int not null default 1,
   is_delete            tinyint not null default 2,
   primary key (id)
);

/*==============================================================*/
/* Index: index_ppm_prs_issue_object_type_org_id                */
/*==============================================================*/
create index index_ppm_prs_issue_object_type_org_id on ppm_prs_issue_object_type
(
   org_id
);

/*===============================================================*/
/* Index: index_ppm_prs_issue_object_type_project_object_type_id */
/*===============================================================*/
create index index_ppm_prs_issue_object_type_project_object_type_id on ppm_prs_issue_object_type
(
   project_object_type_id
);

/*==============================================================*/
/* Index: index_ppm_prs_issue_object_type_create_time           */
/*==============================================================*/
create index index_ppm_prs_issue_object_type_create_time on ppm_prs_issue_object_type
(
   create_time
);

/*==============================================================*/
/* Table: ppm_prs_issue_source                                  */
/*==============================================================*/
create table if not exists ppm_prs_issue_source
(
   id                   bigint not null,
   org_id               bigint not null default 0,
   lang_code            varchar(64) not null default '',
   name                 varchar(64) not null default '',
   sort                 int not null default 0,
   project_object_type_id bigint not null default 0,
   remark               varchar(512) not null default '',
   status               tinyint not null default 1,
   creator              bigint not null default 0,
   create_time          datetime not null default CURRENT_TIMESTAMP,
   updator              bigint not null default 0,
   update_time          datetime not null default CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
   version              int not null default 1,
   is_delete            tinyint not null default 2,
   primary key (id)
);

/*==============================================================*/
/* Index: index_ppm_prs_issue_source_org_id                     */
/*==============================================================*/
create index index_ppm_prs_issue_source_org_id on ppm_prs_issue_source
(
   org_id
);

/*==============================================================*/
/* Index: index_ppm_prs_issue_source_project_object_type_id     */
/*==============================================================*/
create index index_ppm_prs_issue_source_project_object_type_id on ppm_prs_issue_source
(
   project_object_type_id
);

/*==============================================================*/
/* Index: index_ppm_prs_issue_source_create_time                */
/*==============================================================*/
create index index_ppm_prs_issue_source_create_time on ppm_prs_issue_source
(
   create_time
);

/*==============================================================*/
/* Table: ppm_prs_priority                                      */
/*==============================================================*/
create table if not exists ppm_prs_priority
(
   id                   bigint not null,
   org_id               bigint not null default 0,
   lang_code            varchar(64) not null default '',
   name                 varchar(32) not null default '',
   type                 tinyint not null default 1,
   sort                 int not null default 1,
   bg_style             varchar(8) not null default '',
   font_style           varchar(8) not null default '',
   is_default           tinyint not null default 2,
   remark               varchar(512) not null default '',
   status               tinyint not null default 1,
   creator              bigint not null default 0,
   create_time          datetime not null default CURRENT_TIMESTAMP,
   updator              bigint not null default 0,
   update_time          datetime not null default CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
   version              int not null default 1,
   is_delete            tinyint not null default 2,
   primary key (id)
);

/*==============================================================*/
/* Index: index_ppm_prs_priority_org_id                         */
/*==============================================================*/
create index index_ppm_prs_priority_org_id on ppm_prs_priority
(
   org_id
);

/*==============================================================*/
/* Index: index_ppm_prs_priority_create_time                    */
/*==============================================================*/
create index index_ppm_prs_priority_create_time on ppm_prs_priority
(
   create_time
);

/*==============================================================*/
/* Table: ppm_prs_process                                       */
/*==============================================================*/
create table if not exists ppm_prs_process
(
   id                   bigint not null,
   org_id               bigint not null default 0,
   lang_code            varchar(64) not null default '',
   name                 varchar(64) not null default '',
   is_default           tinyint not null default 0,
   type                 tinyint not null default 1,
   sort                 int not null default 1,
   creator              bigint not null default 0,
   create_time          datetime not null default CURRENT_TIMESTAMP,
   updator              bigint not null default 0,
   update_time          datetime not null default CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
   version              int not null default 1,
   is_delete            tinyint not null default 2,
   primary key (id)
);

/*==============================================================*/
/* Index: index_ppm_prs_process_org_id                          */
/*==============================================================*/
create index index_ppm_prs_process_org_id on ppm_prs_process
(
   org_id,
   name
);

/*==============================================================*/
/* Index: index_ppm_prs_process_create_time                     */
/*==============================================================*/
create index index_ppm_prs_process_create_time on ppm_prs_process
(
   create_time
);

/*==============================================================*/
/* Table: ppm_prs_process_process_status                        */
/*==============================================================*/
create table if not exists ppm_prs_process_process_status
(
   id                   bigint not null,
   org_id               bigint not null default 0,
   process_id           bigint not null default 0,
   process_status_id    bigint not null default 0,
   is_init_status       tinyint not null default 2,
   sort                 int not null default 1,
   creator              bigint not null default 0,
   create_time          datetime not null default CURRENT_TIMESTAMP,
   updator              bigint not null default 0,
   update_time          datetime not null default CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
   version              int not null default 1,
   is_delete            tinyint not null default 2,
   primary key (id)
);

/*==============================================================*/
/* Index: index_ppm_prs_process_process_status_org_id           */
/*==============================================================*/
create index index_ppm_prs_process_process_status_org_id on ppm_prs_process_process_status
(
   org_id
);

/*==============================================================*/
/* Index: index_ppm_prs_process_process_status_process_id       */
/*==============================================================*/
create index index_ppm_prs_process_process_status_process_id on ppm_prs_process_process_status
(
   process_id
);

/*==============================================================*/
/* Index: index_ppm_prs_process_process_status_create_time      */
/*==============================================================*/
create index index_ppm_prs_process_process_status_create_time on ppm_prs_process_process_status
(
   create_time
);

/*==============================================================*/
/* Table: ppm_prs_process_status                                */
/*==============================================================*/
create table if not exists ppm_prs_process_status
(
   id                   bigint not null,
   org_id               bigint not null default 0,
   lang_code            varchar(64) not null default '',
   name                 varchar(32) not null default '',
   sort                 int not null default 1,
   bg_style             varchar(8) not null default '',
   font_style           varchar(8) not null default '',
   type                 tinyint not null default 1,
   category             tinyint not null default 1,
   remark               varchar(512) not null default '',
   status               tinyint not null default 1,
   creator              bigint not null default 0,
   create_time          datetime not null default CURRENT_TIMESTAMP,
   updator              bigint not null default 0,
   update_time          datetime not null default CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
   version              int not null default 1,
   is_delete            tinyint not null default 2,
   primary key (id)
);

/*==============================================================*/
/* Index: index_ppm_prs_process_status_org_id                   */
/*==============================================================*/
create index index_ppm_prs_process_status_org_id on ppm_prs_process_status
(
   org_id
);

/*==============================================================*/
/* Index: index_ppm_prs_process_status_create_time              */
/*==============================================================*/
create index index_ppm_prs_process_status_create_time on ppm_prs_process_status
(
   create_time
);

/*==============================================================*/
/* Table: ppm_prs_process_step                                  */
/*==============================================================*/
create table if not exists ppm_prs_process_step
(
   id                   bigint not null,
   org_id               bigint not null default 0,
   process_id           bigint not null default 0,
   lang_code            varchar(64) not null default '',
   name                 varchar(32) not null default '',
   start_status         bigint not null default 0,
   end_status           bigint not null default 0,
   sort                 int not null default 1,
   is_default           tinyint not null default 2,
   remark               varchar(512) not null default '',
   status               tinyint not null default 1,
   creator              bigint not null default 0,
   create_time          datetime not null default CURRENT_TIMESTAMP,
   updator              bigint not null default 0,
   update_time          datetime not null default CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
   version              int not null default 1,
   is_delete            tinyint not null default 2,
   primary key (id)
);

/*==============================================================*/
/* Index: index_ppm_prs_process_step_org_id                     */
/*==============================================================*/
create index index_ppm_prs_process_step_org_id on ppm_prs_process_step
(
   org_id
);

/*==============================================================*/
/* Index: index_ppm_prs_process_step_process_id                 */
/*==============================================================*/
create index index_ppm_prs_process_step_process_id on ppm_prs_process_step
(
   process_id
);

/*==============================================================*/
/* Index: index_ppm_prs_process_step_start_status               */
/*==============================================================*/
create index index_ppm_prs_process_step_start_status on ppm_prs_process_step
(
   start_status
);

/*==============================================================*/
/* Index: index_ppm_prs_process_step_end_status                 */
/*==============================================================*/
create index index_ppm_prs_process_step_end_status on ppm_prs_process_step
(
   end_status
);

/*==============================================================*/
/* Index: index_ppm_prs_process_step_create_time                */
/*==============================================================*/
create index index_ppm_prs_process_step_create_time on ppm_prs_process_step
(
   create_time
);

/*==============================================================*/
/* Table: ppm_prs_project_object_type                           */
/*==============================================================*/
create table if not exists ppm_prs_project_object_type
(
   id                   bigint not null,
   org_id               bigint not null default 0,
   lang_code            varchar(64) not null default '',
   pre_code             varchar(16) not null default '',
   name                 varchar(32) not null default '',
   icon                 varchar(8) not null default '',
   bg_style             varchar(8) not null default '',
   font_style           varchar(8) not null default '',
   object_type          tinyint not null default 1,
   sort                 int not null default 1,
   remark               varchar(512) not null default '',
   is_readonly          tinyint not null default 2,
   status               tinyint not null default 1,
   creator              bigint not null default 0,
   create_time          datetime not null default CURRENT_TIMESTAMP,
   updator              bigint not null default 0,
   update_time          datetime not null default CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
   version              int not null default 1,
   is_delete            tinyint not null default 2,
   primary key (id)
);

/*==============================================================*/
/* Index: index_ppm_prs_project_object_type_org_id_pre_code     */
/*==============================================================*/
create index index_ppm_prs_project_object_type_org_id_pre_code on ppm_prs_project_object_type
(
   org_id,
   pre_code
);

/*==============================================================*/
/* Index: index_ppm_prs_project_objct_type_create_time          */
/*==============================================================*/
create index index_ppm_prs_project_objct_type_create_time on ppm_prs_project_object_type
(
   create_time
);

/*==============================================================*/
/* Table: ppm_prs_project_object_type_process                   */
/*==============================================================*/
create table if not exists ppm_prs_project_object_type_process
(
   id                   bigint not null,
   org_id               bigint not null default 0,
   project_id           bigint not null default 0,
   project_object_type_id bigint not null default 0,
   process_id           bigint not null default 0,
   sort                 int not null default 1,
   creator              bigint not null default 0,
   create_time          datetime not null default CURRENT_TIMESTAMP,
   updator              bigint not null default 0,
   update_time          datetime not null default CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
   version              int not null default 1,
   is_delete            tinyint not null default 2,
   primary key (id)
);

/*==============================================================*/
/* Index: index_ppm_prs_project_object_type_process_org_id      */
/*==============================================================*/
create index index_ppm_prs_project_object_type_process_org_id on ppm_prs_project_object_type_process
(
   org_id,
   project_id
);

/*==============================================================*/
/* Index: index_ppm_prs_project_object_type_process_project_id  */
/*==============================================================*/
create index index_ppm_prs_project_object_type_process_project_id on ppm_prs_project_object_type_process
(
   project_id,
   project_object_type_id
);

/*==============================================================*/
/* Index: index_ppm_prs_project_object_type_process_create_time */
/*==============================================================*/
create index index_ppm_prs_project_object_type_process_create_time on ppm_prs_project_object_type_process
(
   create_time
);

/*==============================================================*/
/* Table: ppm_prs_project_type                                  */
/*==============================================================*/
create table if not exists ppm_prs_project_type
(
   id                   bigint not null,
   org_id               bigint not null default 0,
   lang_code            varchar(64) not null default '',
   name                 varchar(32) not null default '',
   sort                 int not null default 1,
   default_process_id   bigint not null default 0,
   category             bigint not null default 0,
   mode                 tinyint not null default 1,
   is_readonly          tinyint not null default 2,
   remark               varchar(512) not null default '',
   status               tinyint not null default 1,
   creator              bigint not null default 0,
   create_time          datetime not null default CURRENT_TIMESTAMP,
   updator              bigint not null default 0,
   update_time          datetime not null default CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
   version              int not null default 1,
   is_delete            tinyint not null default 2
);

/*==============================================================*/
/* Index: index_ppm_prs_project_type_org_id                     */
/*==============================================================*/
create index index_ppm_prs_project_type_org_id on ppm_prs_project_type
(
   org_id
);

/*==============================================================*/
/* Index: index_ppm_prs_project_type_create_time                */
/*==============================================================*/
create index index_ppm_prs_project_type_create_time on ppm_prs_project_type
(
   create_time
);

/*==============================================================*/
/* Table: ppm_prs_project_type_project_object_type              */
/*==============================================================*/
create table if not exists ppm_prs_project_type_project_object_type
(
   id                   bigint not null,
   org_id               bigint not null default 0,
   project_type_id      bigint not null default 0,
   project_object_type_id bigint not null default 0,
   remark               varchar(512) not null default '',
   default_process_id   bigint not null default 0,
   is_readonly          tinyint not null default 2,
   status               tinyint not null default 1,
   creator              bigint not null default 0,
   create_time          datetime not null default CURRENT_TIMESTAMP,
   updator              bigint not null default 0,
   update_time          datetime not null default CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
   version              int not null default 1,
   is_delete            tinyint not null default 2,
   primary key (id)
);

/*=========================================================================*/
/* Index: index_ppm_prs_project_type_project_object_type_org_id_project_ty */
/*=========================================================================*/
create index index_ppm_prs_project_type_project_object_type_org_id_project_ty on ppm_prs_project_type_project_object_type
(
   org_id,
   project_type_id
);

/*==================================================================*/
/* Index: index_ppm_prs_project_type_project_objct_type_create_time */
/*==================================================================*/
create index index_ppm_prs_project_type_project_objct_type_create_time on ppm_prs_project_type_project_object_type
(
   create_time
);

/*==============================================================*/
/* Table: ppm_res_resource                                      */
/*==============================================================*/
create table if not exists ppm_res_resource
(
   id                   bigint not null,
   org_id               bigint not null default 0,
   project_id           bigint not null default 0,
   type                 tinyint(2) not null default 1,
   bucket               varchar(32) not null default '',
   path                 varchar(512) not null default '1',
   name                 varchar(256) not null default '',
   suffix               varchar(32) not null default 'txt',
   md5                  varchar(32) not null default '',
   size                 bigint not null default 0,
   creator              bigint not null default 0,
   create_time          datetime not null default CURRENT_TIMESTAMP,
   updator              bigint not null default 0,
   update_time          datetime not null default CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
   version              int not null default 1,
   is_delete            tinyint not null default 2,
   primary key (id)
);

/*==============================================================*/
/* Index: index_ppm_res_resource_md5                            */
/*==============================================================*/
create index index_ppm_res_resource_md5 on ppm_res_resource
(
   md5
);

/*==============================================================*/
/* Index: index_ppm_res_resource_org_id                         */
/*==============================================================*/
create index index_ppm_res_resource_org_id on ppm_res_resource
(
   org_id
);

/*==============================================================*/
/* Index: index_ppm_res_resource_project_id                     */
/*==============================================================*/
create index index_ppm_res_resource_project_id on ppm_res_resource
(
   project_id
);

/*==============================================================*/
/* Index: index_ppm_res_resource_create_time                    */
/*==============================================================*/
create index index_ppm_res_resource_create_time on ppm_res_resource
(
   create_time
);

/*==============================================================*/
/* Index: index_ppm_res_resource_creator                        */
/*==============================================================*/
create index index_ppm_res_resource_creator on ppm_res_resource
(
   creator
);

/*==============================================================*/
/* Table: ppm_rol_operation                                     */
/*==============================================================*/
create table if not exists ppm_rol_operation
(
   id                   bigint not null,
   code                 varchar(16) not null default '',
   name                 varchar(128) not null default '',
   remark               varchar(512) not null default '',
   status               tinyint not null default 1,
   creator              bigint not null default 0,
   create_time          datetime not null default CURRENT_TIMESTAMP,
   updator              bigint not null default 0,
   update_time          datetime not null default CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
   version              int not null default 1,
   is_delete            tinyint not null default 2,
   primary key (id)
);

/*==============================================================*/
/* Index: index_ppm_rol_operation_code                          */
/*==============================================================*/
create index index_ppm_rol_operation_code on ppm_rol_operation
(
   code
);

/*==============================================================*/
/* Index: index_ppm_rol_operation_create_time                   */
/*==============================================================*/
create index index_ppm_rol_operation_create_time on ppm_rol_operation
(
   create_time
);

/*==============================================================*/
/* Table: ppm_rol_permission                                    */
/*==============================================================*/
create table if not exists ppm_rol_permission
(
   id                   bigint not null,
   org_id               bigint not null default 0,
   lang_code            varchar(64) not null default '',
   code                 varchar(32) not null default '',
   name                 varchar(128) not null default '',
   parent_id            bigint not null default 0,
   type                 tinyint not null default 1,
   path                 varchar(512) not null default '',
   is_show              tinyint not null default 1,
   remark               varchar(512) not null default '',
   status               tinyint not null default 1,
   creator              bigint not null default 0,
   create_time          datetime not null default CURRENT_TIMESTAMP,
   updator              bigint not null default 0,
   update_time          datetime not null default CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
   version              int not null default 1,
   is_delete            tinyint not null default 2,
   primary key (id)
);

/*==============================================================*/
/* Index: index_ppm_rol_permission_org_id                       */
/*==============================================================*/
create index index_ppm_rol_permission_org_id on ppm_rol_permission
(
   org_id
);

/*==============================================================*/
/* Index: index_ppm_rol_permission_parent_id                    */
/*==============================================================*/
create index index_ppm_rol_permission_parent_id on ppm_rol_permission
(
   parent_id
);

/*==============================================================*/
/* Index: index_ppm_rol_permission_path                         */
/*==============================================================*/
create index index_ppm_rol_permission_path on ppm_rol_permission
(
   path
);

/*==============================================================*/
/* Index: index_ppm_rol_permission_create_time                  */
/*==============================================================*/
create index index_ppm_rol_permission_create_time on ppm_rol_permission
(
   create_time
);

/*==============================================================*/
/* Table: ppm_rol_permission_operation                          */
/*==============================================================*/
create table if not exists ppm_rol_permission_operation
(
   id                   bigint not null,
   org_id               bigint not null default 0,
   permission_id        bigint not null default 0,
   lang_code            varchar(64) not null default '',
   name                 varchar(128) not null default '',
   operation_codes      varchar(128) not null default '',
   remark               varchar(512) not null default '',
   is_show              tinyint not null default 1,
   status               tinyint not null default 1,
   creator              bigint not null default 0,
   create_time          datetime not null default CURRENT_TIMESTAMP,
   updator              bigint not null default 0,
   update_time          datetime not null default CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
   version              int not null default 1,
   is_delete            tinyint not null default 2,
   primary key (id)
);

/*==============================================================*/
/* Index: index_ppm_rol_permission_operation_org_id             */
/*==============================================================*/
create index index_ppm_rol_permission_operation_org_id on ppm_rol_permission_operation
(
   org_id
);

/*==============================================================*/
/* Index: index_ppm_rol_permission_operation_permission_id      */
/*==============================================================*/
create index index_ppm_rol_permission_operation_permission_id on ppm_rol_permission_operation
(
   permission_id
);

/*==============================================================*/
/* Index: index_ppm_rol_permission_operation_create_time        */
/*==============================================================*/
create index index_ppm_rol_permission_operation_create_time on ppm_rol_permission_operation
(
   create_time
);

/*==============================================================*/
/* Table: ppm_rol_role                                          */
/*==============================================================*/
create table if not exists ppm_rol_role
(
   id                   bigint not null,
   org_id               bigint not null default 0,
   lang_code            varchar(64) not null default '',
   name                 varchar(64) not null default '',
   remark               varchar(512) not null default '',
   is_readonly          tinyint not null default 2,
   is_modify_permission tinyint not null default 1,
   is_default           tinyint not null default 2,
   role_group_id        bigint not null default 0,
   status               tinyint not null default 1,
   creator              bigint not null default 0,
   create_time          datetime not null default CURRENT_TIMESTAMP,
   updator              bigint not null default 0,
   update_time          datetime not null default CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
   version              int not null default 1,
   is_delete            tinyint not null default 2,
   primary key (id)
);

/*==============================================================*/
/* Index: index_ppm_rol_role_org_id                             */
/*==============================================================*/
create index index_ppm_rol_role_org_id on ppm_rol_role
(
   org_id
);

/*==============================================================*/
/* Index: index_ppm_rol_role_role_group_id                      */
/*==============================================================*/
create index index_ppm_rol_role_role_group_id on ppm_rol_role
(
   role_group_id
);

/*==============================================================*/
/* Index: index_ppm_rol_role_create_time                        */
/*==============================================================*/
create index index_ppm_rol_role_create_time on ppm_rol_role
(
   create_time
);

/*==============================================================*/
/* Table: ppm_rol_role_group                                    */
/*==============================================================*/
create table if not exists ppm_rol_role_group
(
   id                   bigint not null,
   org_id               bigint not null default 0,
   lang_code            varchar(64) not null default '',
   name                 varchar(64) not null default '',
   remark               varchar(512) not null default '',
   type                 tinyint not null default 1,
   is_readonly          tinyint not null default 2,
   is_show              tinyint not null default 1,
   is_default           tinyint not null default 2,
   status               tinyint not null default 1,
   creator              bigint not null default 0,
   create_time          datetime not null default CURRENT_TIMESTAMP,
   updator              bigint not null default 0,
   update_time          datetime not null default CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
   version              int not null default 1,
   is_delete            tinyint not null default 2,
   primary key (id)
);

/*==============================================================*/
/* Index: index_ppm_rol_role_group_org_id                       */
/*==============================================================*/
create index index_ppm_rol_role_group_org_id on ppm_rol_role_group
(
   org_id
);

/*==============================================================*/
/* Index: index_ppm_rol_role_group_create_time                  */
/*==============================================================*/
create index index_ppm_rol_role_group_create_time on ppm_rol_role_group
(
   create_time
);

/*==============================================================*/
/* Table: ppm_rol_role_permission_operation                     */
/*==============================================================*/
create table if not exists ppm_rol_role_permission_operation
(
   id                   bigint not null,
   org_id               bigint not null default 0,
   role_id              bigint not null default 0,
   project_id           bigint not null default 0,
   permission_id        bigint not null default 0,
   permission_path      varchar(512) not null default '',
   operation_codes      varchar(128) not null default '',
   creator              bigint not null default 0,
   create_time          datetime not null default CURRENT_TIMESTAMP,
   updator              bigint not null default 0,
   update_time          datetime not null default CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
   version              int not null default 1,
   is_delete            tinyint not null default 2,
   primary key (id)
);

/*==============================================================*/
/* Index: index_ppm_rol_role_permission_operation_org_id        */
/*==============================================================*/
create index index_ppm_rol_role_permission_operation_org_id on ppm_rol_role_permission_operation
(
   org_id,
   project_id
);

/*==================================================================*/
/* Index: index_ppm_rol_role_permission_operation_org_permission_id */
/*==================================================================*/
create index index_ppm_rol_role_permission_operation_org_permission_id on ppm_rol_role_permission_operation
(
   permission_id
);

/*====================================================================*/
/* Index: index_ppm_rol_role_permission_operation_org_permission_path */
/*====================================================================*/
create index index_ppm_rol_role_permission_operation_org_permission_path on ppm_rol_role_permission_operation
(
   permission_path
);

/*================================================================*/
/* Index: index_ppm_rol_role_permission_operation_org_create_time */
/*================================================================*/
create index index_ppm_rol_role_permission_operation_org_create_time on ppm_rol_role_permission_operation
(
   create_time
);

/*==============================================================*/
/* Table: ppm_rol_role_user                                     */
/*==============================================================*/
create table if not exists ppm_rol_role_user
(
   id                   bigint not null,
   org_id               bigint not null default 0,
   project_id           bigint not null default 0,
   role_id              bigint not null default 0,
   user_id              bigint not null default 0,
   creator              bigint not null default 0,
   create_time          datetime not null default CURRENT_TIMESTAMP,
   updator              bigint not null default 0,
   update_time          datetime not null default CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
   version              int not null default 1,
   is_delete            tinyint not null default 2,
   primary key (id)
);

/*==============================================================*/
/* Index: index_ppm_rol_role_user_org_id                        */
/*==============================================================*/
create index index_ppm_rol_role_user_org_id on ppm_rol_role_user
(
   org_id
);

/*==============================================================*/
/* Index: index_ppm_rol_role_user_role_id                       */
/*==============================================================*/
create index index_ppm_rol_role_user_role_id on ppm_rol_role_user
(
   role_id
);

/*==============================================================*/
/* Index: index_ppm_rol_role_user_user_id_project_id            */
/*==============================================================*/
create index index_ppm_rol_role_user_user_id_project_id on ppm_rol_role_user
(
   user_id,
   project_id
);

/*==============================================================*/
/* Index: index_ppm_rol_role_user_create_time                   */
/*==============================================================*/
create index index_ppm_rol_role_user_create_time on ppm_rol_role_user
(
   create_time
);

/*==============================================================*/
/* Table: ppm_sha_share                                         */
/*==============================================================*/
create table if not exists ppm_sha_share
(
   id                   bigint not null,
   org_id               bigint not null default 0,
   project_id           bigint not null default 0,
   name                 varchar(128) not null default '',
   remark               varchar(512) not null default '',
   logo                 varchar(512) not null default '',
   type                 int not null default 1,
   content              text,
   content_md5          varchar(32) not null default '',
   finish_time          datetime not null default '1970-01-01 00:00:00',
   status               tinyint not null default 1,
   creator              bigint not null default 0,
   create_time          datetime not null default CURRENT_TIMESTAMP,
   updator              bigint not null default 0,
   update_time          datetime not null default CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
   version              int not null default 1,
   is_delete            tinyint not null default 2,
   primary key (id)
);

/*==============================================================*/
/* Index: index_ppm_sha_share_org_id                            */
/*==============================================================*/
create index index_ppm_sha_share_org_id on ppm_sha_share
(
   org_id
);

/*==============================================================*/
/* Index: index_ppm_sha_share_project_id                        */
/*==============================================================*/
create index index_ppm_sha_share_project_id on ppm_sha_share
(
   project_id
);

/*==============================================================*/
/* Index: index_ppm_sha_share_finish_time                       */
/*==============================================================*/
create index index_ppm_sha_share_finish_time on ppm_sha_share
(
   finish_time
);

/*==============================================================*/
/* Index: index_ppm_sha_share_content_md5                       */
/*==============================================================*/
create index index_ppm_sha_share_content_md5 on ppm_sha_share
(
   content_md5
);

/*==============================================================*/
/* Index: index_ppm_sha_share_creator                           */
/*==============================================================*/
create index index_ppm_sha_share_creator on ppm_sha_share
(
   creator
);

/*==============================================================*/
/* Index: index_ppm_sha_share_create_time                       */
/*==============================================================*/
create index index_ppm_sha_share_create_time on ppm_sha_share
(
   create_time
);

/*==============================================================*/
/* Table: ppm_sta_iteration_stat                                */
/*==============================================================*/
create table if not exists ppm_sta_iteration_stat
(
   id                   bigint not null,
   org_id               bigint not null default 0,
   project_id           bigint not null default 0,
   iteration_id         bigint not null default 0,
   issue_count          int not null default 0,
   issue_wait_count     int not null default 0,
   issue_running_count  int not null default 0,
   issue_end_count      int not null default 0,
   demand_count         int not null default 0,
   demand_wait_count    int not null default 0,
   demand_running_count int not null default 0,
   demand_end_count     int not null default 0,
   story_point_count    int not null default 0,
   story_point_wait_count int not null default 0,
   story_point_running_count int not null default 0,
   story_point_end_count int not null default 0,
   task_count           int not null default 0,
   task_wait_count      int not null default 0,
   task_running_count   int not null default 0,
   task_end_count       int not null default 0,
   bug_count            int not null default 0,
   bug_wait_count       int not null default 0,
   bug_running_count    int not null default 0,
   bug_end_count        int not null default 0,
   testtask_count       int not null default 0,
   testtask_wait_count  int not null default 0,
   testtask_running_count int not null default 0,
   testtask_end_count   int not null default 0,
   ext                  varchar(4096) not null default '',
   stat_date            date not null default '1970-01-01',
   status               bigint not null default 0,
   creator              bigint not null default 0,
   create_time          datetime not null default CURRENT_TIMESTAMP,
   updator              bigint not null default 0,
   update_time          datetime not null default CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
   version              int not null default 1,
   is_delete            tinyint not null default 2,
   primary key (id)
);

/*==============================================================*/
/* Index: index_ppm_sta_iteration_stat_org_id                   */
/*==============================================================*/
create index index_ppm_sta_iteration_stat_org_id on ppm_sta_iteration_stat
(
   org_id
);

/*==============================================================*/
/* Index: index_ppm_sta_iteration_stat_iteration_id             */
/*==============================================================*/
create index index_ppm_sta_iteration_stat_iteration_id on ppm_sta_iteration_stat
(
   iteration_id
);

/*==============================================================*/
/* Index: index_ppm_sta_iteration_stat_project_id               */
/*==============================================================*/
create index index_ppm_sta_iteration_stat_project_id on ppm_sta_iteration_stat
(
   project_id
);

/*==============================================================*/
/* Index: index_ppm_sta_iteration_stat_create_time              */
/*==============================================================*/
create index index_ppm_sta_iteration_stat_create_time on ppm_sta_iteration_stat
(
   create_time
);

/*==============================================================*/
/* Index: index_ppm_sta_iteration_stat_stat_date                */
/*==============================================================*/
create index index_ppm_sta_iteration_stat_stat_date on ppm_sta_iteration_stat
(
   stat_date
);

/*==============================================================*/
/* Table: ppm_sta_project_day_stat                              */
/*==============================================================*/
create table if not exists ppm_sta_project_day_stat
(
   id                   bigint not null,
   org_id               bigint not null default 0,
   project_id           bigint not null default 0,
   issue_count          int not null default 0,
   issue_wait_count     int not null default 0,
   issue_running_count  int not null default 0,
   issue_overdue_count  int not null default 0,
   issue_end_count      int not null default 0,
   demand_count         int not null default 0,
   demand_wait_count    int not null default 0,
   demand_running_count int not null default 0,
   demand_overdue_count int not null default 0,
   demand_end_count     int not null default 0,
   story_point_count    int not null default 0,
   story_point_wait_count int not null default 0,
   story_point_running_count int not null default 0,
   story_point_overdue_count int not null default 0,
   story_point_end_count int not null default 0,
   task_count           int not null default 0,
   task_wait_count      int not null default 0,
   task_running_count   int not null default 0,
   task_overdue_count   int not null default 0,
   task_end_count       int not null default 0,
   bug_count            int not null default 0,
   bug_wait_count       int not null default 0,
   bug_running_count    int not null default 0,
   bug_overdue_count    int not null default 0,
   bug_end_count        int not null default 0,
   testtask_count       int not null default 0,
   testtask_wait_count  int not null default 0,
   testtask_running_count int not null default 0,
   testtask_overdue_count int not null default 0,
   testtask_end_count   int not null default 0,
   ext                  varchar(4096) not null default '',
   stat_date            date not null default '1970-01-01',
   status               bigint not null default 0,
   creator              bigint not null default 0,
   create_time          datetime not null default CURRENT_TIMESTAMP,
   updator              bigint not null default 0,
   update_time          datetime not null default CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
   version              int not null default 1,
   is_delete            tinyint not null default 2,
   primary key (id)
);

/*==============================================================*/
/* Index: ppm_sta_project_day_stat_org_id                       */
/*==============================================================*/
create index ppm_sta_project_day_stat_org_id on ppm_sta_project_day_stat
(
   org_id
);

/*==============================================================*/
/* Index: ppm_sta_project_day_stat_project_id                   */
/*==============================================================*/
create index ppm_sta_project_day_stat_project_id on ppm_sta_project_day_stat
(
   project_id
);

/*==============================================================*/
/* Index: ppm_sta_project_day_stat_create_time                  */
/*==============================================================*/
create index ppm_sta_project_day_stat_create_time on ppm_sta_project_day_stat
(
   create_time
);

/*==============================================================*/
/* Index: index_ppm_sta_project_day_stat_stat_date              */
/*==============================================================*/
create index index_ppm_sta_project_day_stat_stat_date on ppm_sta_project_day_stat
(
   stat_date
);

/*==============================================================*/
/* Table: ppm_tak_message                                       */
/*==============================================================*/
create table if not exists ppm_tak_message
(
   id                   bigint not null,
   org_id               bigint not null default 0,
   type                 int not null default 0,
   project_id           bigint not null default 0,
   issue_id             bigint not null default 0,
   trends_id            bigint not null default 0,
   info                 varchar(4096) not null default '',
   content              text,
   fail_count           int not null default 0,
   fail_time            datetime not null default '1970-01-01 00:00:00',
   fail_msg             varchar(1024) not null default '',
   finish_status        varchar(32) not null default '',
   finish_msg           varchar(1024) not null default '',
   start_time           datetime not null default '1970-01-01 00:00:00',
   status               tinyint not null default 1,
   creator              bigint not null default 0,
   create_time          datetime not null default CURRENT_TIMESTAMP,
   updator              bigint not null default 0,
   update_time          datetime not null default CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
   version              int not null default 1,
   is_delete            tinyint not null default 2,
   primary key (id)
);

/*==============================================================*/
/* Index: index_ppm_tak_message_issue_id                        */
/*==============================================================*/
create index index_ppm_tak_message_issue_id on ppm_tak_message
(
   issue_id
);

/*==============================================================*/
/* Index: index_ppm_tak_message_project_id                      */
/*==============================================================*/
create index index_ppm_tak_message_project_id on ppm_tak_message
(
   project_id
);

/*==============================================================*/
/* Index: index_ppm_tak_message_org_id                          */
/*==============================================================*/
create index index_ppm_tak_message_org_id on ppm_tak_message
(
   org_id
);

/*==============================================================*/
/* Index: index_ppm_tak_message_create_time                     */
/*==============================================================*/
create index index_ppm_tak_message_create_time on ppm_tak_message
(
   create_time
);

/*==============================================================*/
/* Index: index_ppm_tak_message_trends_id                       */
/*==============================================================*/
create index index_ppm_tak_message_trends_id on ppm_tak_message
(
   trends_id
);

/*==============================================================*/
/* Table: ppm_tem_team                                          */
/*==============================================================*/
create table if not exists ppm_tem_team
(
   id                   bigint not null,
   org_id               bigint not null default 0,
   name                 varchar(128) not null default '',
   nick_name            varchar(128) not null default '',
   owner                bigint not null default 0,
   department_id        bigint not null default 0,
   is_default           tinyint not null default 2,
   remark               varchar(512) not null default '',
   status               tinyint not null default 1,
   creator              bigint not null default 0,
   create_time          datetime not null default CURRENT_TIMESTAMP,
   updator              bigint not null default 0,
   update_time          datetime not null default CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
   version              int not null default 1,
   is_delete            tinyint not null default 2,
   primary key (id)
);

/*==============================================================*/
/* Index: index_ppm_tem_team_org_id                             */
/*==============================================================*/
create index index_ppm_tem_team_org_id on ppm_tem_team
(
   org_id
);

/*==============================================================*/
/* Index: index_ppm_tem_team_create_time                        */
/*==============================================================*/
create index index_ppm_tem_team_create_time on ppm_tem_team
(
   create_time
);

/*==============================================================*/
/* Index: index_ppm_tem_team_department_id                      */
/*==============================================================*/
create index index_ppm_tem_team_department_id on ppm_tem_team
(
   department_id
);

/*==============================================================*/
/* Table: ppm_tem_user_team                                     */
/*==============================================================*/
create table if not exists ppm_tem_user_team
(
   id                   bigint not null,
   org_id               bigint not null default 0,
   team_id              bigint not null default 0,
   user_id              bigint not null default 0,
   relation_type        tinyint not null default 1,
   status               tinyint not null default 1,
   creator              bigint not null default 0,
   create_time          datetime not null default CURRENT_TIMESTAMP,
   updator              bigint not null default 0,
   update_time          datetime not null default CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
   version              int not null default 1,
   is_delete            tinyint not null default 2,
   primary key (id)
);

/*==============================================================*/
/* Index: index_ppm_tem_user_team_org_id                        */
/*==============================================================*/
create index index_ppm_tem_user_team_org_id on ppm_tem_user_team
(
   org_id
);

/*==============================================================*/
/* Index: index_ppm_tem_user_team_team_id                       */
/*==============================================================*/
create index index_ppm_tem_user_team_team_id on ppm_tem_user_team
(
   team_id
);

/*==============================================================*/
/* Index: index_ppm_tem_user_team_user_id                       */
/*==============================================================*/
create index index_ppm_tem_user_team_user_id on ppm_tem_user_team
(
   user_id
);

/*==============================================================*/
/* Index: index_ppm_tem_user_team_create_time                   */
/*==============================================================*/
create index index_ppm_tem_user_team_create_time on ppm_tem_user_team
(
   create_time
);

/*==============================================================*/
/* Table: ppm_tre_comment                                       */
/*==============================================================*/
create table if not exists ppm_tre_comment
(
   id                   bigint not null,
   org_id               bigint not null default 0,
   project_id           bigint not null default 0,
   trends_id            bigint not null default 0,
   object_id            bigint not null default 0,
   object_type          varchar(32) not null default 'issue',
   content              varchar(4096) not null default '',
   parent_id            bigint not null default 0,
   creator              bigint not null default 0,
   create_time          datetime not null default CURRENT_TIMESTAMP,
   updator              bigint not null default 0,
   update_time          datetime not null default CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
   version              int not null default 1,
   is_delete            tinyint not null default 2,
   primary key (id)
);

/*==============================================================*/
/* Index: index_ppm_tre_comment_org_id                          */
/*==============================================================*/
create index index_ppm_tre_comment_org_id on ppm_tre_comment
(
   org_id
);

/*==============================================================*/
/* Index: index_ppm_tre_comment_create_time                     */
/*==============================================================*/
create index index_ppm_tre_comment_create_time on ppm_tre_comment
(
   create_time
);

/*==============================================================*/
/* Index: index_ppm_tre_comment_project_id                      */
/*==============================================================*/
create index index_ppm_tre_comment_project_id on ppm_tre_comment
(
   project_id
);

/*==============================================================*/
/* Index: index_ppm_tre_comment_trends_id                       */
/*==============================================================*/
create index index_ppm_tre_comment_trends_id on ppm_tre_comment
(
   trends_id
);

/*==============================================================*/
/* Index: index_ppm_tre_comment_parent_id                       */
/*==============================================================*/
create index index_ppm_tre_comment_parent_id on ppm_tre_comment
(
   parent_id
);

/*==============================================================*/
/* Index: index_ppm_tre_comment_object_id                       */
/*==============================================================*/
create index index_ppm_tre_comment_object_id on ppm_tre_comment
(
   object_id
);

/*==============================================================*/
/* Index: index_ppm_tre_trends_oper_create_time                 */
/*==============================================================*/
create index index_ppm_tre_trends_oper_create_time on ppm_tre_comment
(
   create_time
);

/*==============================================================*/
/* Index: index_ppm_tre_trends_oper_creator                     */
/*==============================================================*/
create index index_ppm_tre_trends_oper_creator on ppm_tre_comment
(
   creator
);

/*==============================================================*/
/* Table: ppm_tre_trends                                        */
/*==============================================================*/
create table if not exists ppm_tre_trends
(
   id                   bigint not null,
   org_id               bigint not null default 0,
   uuid                 varchar(32) not null default '',
   module1              varchar(16) not null default '',
   module2_id           bigint not null,
   module2              varchar(16) not null default '',
   module3_id           bigint not null default 0,
   module3              varchar(16) not null default '',
   oper_code            varchar(16) not null default '',
   oper_obj_id          bigint not null default 0,
   oper_obj_type        varchar(32) not null default '',
   oper_obj_property    varchar(32) not null default '',
   relation_obj_id      bigint not null default 0,
   relation_obj_type    varchar(32) not null default '',
   relation_type        varchar(32) not null default '',
   new_value            text,
   old_value            text,
   ext                  varchar(4096) not null default '',
   creator              bigint not null default 0,
   create_time          datetime not null default CURRENT_TIMESTAMP,
   is_delete            tinyint not null default 2,
   primary key (id)
);

/*==============================================================*/
/* Index: index_ppm_tre_trends_org_id                           */
/*==============================================================*/
create index index_ppm_tre_trends_org_id on ppm_tre_trends
(
   org_id
);

/*==============================================================*/
/* Index: index_ppm_tre_trends_uuid                             */
/*==============================================================*/
create index index_ppm_tre_trends_uuid on ppm_tre_trends
(
   uuid
);

/*==============================================================*/
/* Index: index_ppm_tre_trends_module2_id_module2               */
/*==============================================================*/
create index index_ppm_tre_trends_module2_id_module2 on ppm_tre_trends
(
   module2_id,
   module2
);

/*==============================================================*/
/* Index: index_ppm_tre_trends_module3_id_module3               */
/*==============================================================*/
create index index_ppm_tre_trends_module3_id_module3 on ppm_tre_trends
(
   module3_id,
   module3
);

/*====================================================================*/
/* Index: index_ppm_tre_trends_oper_relation_obj_id_relation_obj_type */
/*====================================================================*/
create index index_ppm_tre_trends_oper_relation_obj_id_relation_obj_type on ppm_tre_trends
(
   relation_obj_id,
   relation_obj_type
);

/*==============================================================*/
/* Index: index_ppm_tre_trends_oper_obj_id                      */
/*==============================================================*/
create index index_ppm_tre_trends_oper_obj_id on ppm_tre_trends
(
   oper_obj_id,
   oper_obj_type
);

/*==============================================================*/
/* Index: index_ppm_tre_trends_oper_create_time                 */
/*==============================================================*/
create index index_ppm_tre_trends_oper_create_time on ppm_tre_trends
(
   create_time
);

/*==============================================================*/
/* Index: index_ppm_tre_trends_oper_creator                     */
/*==============================================================*/
create index index_ppm_tre_trends_oper_creator on ppm_tre_trends
(
   creator
);

/*==============================================================*/
/* Table: ppm_tst_case                                          */
/*==============================================================*/
create table if not exists ppm_tst_case
(
   id                   bigint not null,
   org_id               bigint not null default 0,
   app_id               bigint not null default 0,
   app_version_id       bigint not null default 0,
   group_id             bigint not null default 0,
   sort                 int not null default 0,
   code                 varchar(128) not null default '0',
   name                 varchar(128) not null default '',
   creator              bigint not null default 0,
   create_time          datetime not null default CURRENT_TIMESTAMP,
   updator              bigint not null default 0,
   update_time          datetime not null default CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
   version              int not null default 1,
   is_delete            tinyint not null default 0,
   primary key (id)
);

/*==============================================================*/
/* Index: index_ppm_tst_case_org_id                             */
/*==============================================================*/
create index index_ppm_tst_case_org_id on ppm_tst_case
(
   org_id
);

/*==============================================================*/
/* Index: index_ppm_tst_case_create_time                        */
/*==============================================================*/
create index index_ppm_tst_case_create_time on ppm_tst_case
(
   create_time
);

/*==============================================================*/
/* Index: index_ppm_tst_pro_id                                  */
/*==============================================================*/
create index index_ppm_tst_pro_id on ppm_tst_case
(
   app_id
);

/*==============================================================*/
/* Index: index_ppm_tst_group_id                                */
/*==============================================================*/
create index index_ppm_tst_group_id on ppm_tst_case
(
   group_id
);

/*==============================================================*/
/* Index: index_ppm_tst_app_version_id                          */
/*==============================================================*/
create index index_ppm_tst_app_version_id on ppm_tst_case
(
   app_version_id
);

/*==============================================================*/
/* Index: index_ppm_tst_sort                                    */
/*==============================================================*/
create index index_ppm_tst_sort on ppm_tst_case
(
   sort
);

/*==============================================================*/
/* Index: index_ppm_tst_case_code                               */
/*==============================================================*/
create index index_ppm_tst_case_code on ppm_tst_case
(
   code
);

/*==============================================================*/
/* Table: ppm_tst_case_attachment                               */
/*==============================================================*/
create table if not exists ppm_tst_case_attachment
(
   id                   bigint not null,
   org_id               bigint not null default 0,
   app_id               bigint not null default 0,
   app_version_id       bigint not null default 0,
   case_id              bigint not null default 0,
   resource_id          bigint not null default 0,
   sort                 int not null default 0,
   creator              bigint not null default 0,
   create_time          datetime not null default CURRENT_TIMESTAMP,
   updator              bigint not null default 0,
   update_time          datetime not null default CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
   version              int not null default 1,
   is_delete            tinyint not null default 2,
   primary key (id)
);

/*==============================================================*/
/* Index: index_ppm_tst_case_org_id                             */
/*==============================================================*/
create index index_ppm_tst_case_org_id on ppm_tst_case_attachment
(
   org_id
);

/*==============================================================*/
/* Index: index_ppm_tst_case_create_time                        */
/*==============================================================*/
create index index_ppm_tst_case_create_time on ppm_tst_case_attachment
(
   create_time
);

/*==============================================================*/
/* Index: index_ppm_tst_case_id                                 */
/*==============================================================*/
create index index_ppm_tst_case_id on ppm_tst_case_attachment
(
   case_id
);

/*==============================================================*/
/* Index: index_ppm_tst_resouce_id                              */
/*==============================================================*/
create index index_ppm_tst_resouce_id on ppm_tst_case_attachment
(
   resource_id
);

/*==============================================================*/
/* Index: index_ppm_tst_pro_id                                  */
/*==============================================================*/
create index index_ppm_tst_pro_id on ppm_tst_case_attachment
(
   app_id
);

/*==============================================================*/
/* Index: index_ppm_tst_app_version_id                          */
/*==============================================================*/
create index index_ppm_tst_app_version_id on ppm_tst_case_attachment
(
   app_version_id
);

/*==============================================================*/
/* Table: ppm_tst_case_detail                                   */
/*==============================================================*/
create table if not exists ppm_tst_case_detail
(
   id                   bigint not null,
   org_id               bigint not null default 0,
   app_id               bigint not null,
   app_version_id       bigint not null default 0,
   case_id              bigint not null,
   pre_condition        text,
   remark               text,
   creator              bigint not null default 0,
   create_time          datetime not null default CURRENT_TIMESTAMP,
   updator              bigint not null default 0,
   update_time          datetime not null default CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
   version              int not null default 1,
   is_delete            tinyint not null default 2,
   primary key (id)
);

/*==============================================================*/
/* Index: index_ppm_tst_case_org_id                             */
/*==============================================================*/
create index index_ppm_tst_case_org_id on ppm_tst_case_detail
(
   org_id
);

/*==============================================================*/
/* Index: index_ppm_tst_case_create_time                        */
/*==============================================================*/
create index index_ppm_tst_case_create_time on ppm_tst_case_detail
(
   create_time
);

/*==============================================================*/
/* Index: index_ppm_tst_case_id                                 */
/*==============================================================*/
create index index_ppm_tst_case_id on ppm_tst_case_detail
(
   case_id
);

/*==============================================================*/
/* Index: index_ppm_tst_pro_id                                  */
/*==============================================================*/
create index index_ppm_tst_pro_id on ppm_tst_case_detail
(
   app_id
);

/*==============================================================*/
/* Index: index_ppm_tst_app_version_id                          */
/*==============================================================*/
create index index_ppm_tst_app_version_id on ppm_tst_case_detail
(
   app_version_id
);

/*==============================================================*/
/* Table: ppm_tst_case_group                                    */
/*==============================================================*/
create table if not exists ppm_tst_case_group
(
   id                   bigint not null,
   org_id               bigint not null default 0,
   app_id               bigint not null default 0,
   app_version_id       bigint not null default 0,
   name                 varchar(128) not null default '',
   parent_id            bigint not null default 0,
   sort                 int not null default 0,
   creator              bigint not null default 0,
   create_time          datetime not null default CURRENT_TIMESTAMP,
   updator              bigint not null default 0,
   update_time          datetime not null default CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
   version              int not null default 1,
   is_delete            tinyint not null default 2,
   primary key (id)
);

/*==============================================================*/
/* Index: index_ppm_tst_case_group_org_id                       */
/*==============================================================*/
create index index_ppm_tst_case_group_org_id on ppm_tst_case_group
(
   org_id
);

/*==============================================================*/
/* Index: index_ppm_tst_case_group_parent_id                    */
/*==============================================================*/
create index index_ppm_tst_case_group_parent_id on ppm_tst_case_group
(
   parent_id
);

/*==============================================================*/
/* Index: index_ppm_tst_case_group_create_time                  */
/*==============================================================*/
create index index_ppm_tst_case_group_create_time on ppm_tst_case_group
(
   create_time
);

/*==============================================================*/
/* Index: index_ppm_tst_case_group_sort                         */
/*==============================================================*/
create index index_ppm_tst_case_group_sort on ppm_tst_case_group
(
   sort
);

/*==============================================================*/
/* Index: index_ppm_tst_project_id                              */
/*==============================================================*/
create index index_ppm_tst_project_id on ppm_tst_case_group
(
   app_id
);

/*==============================================================*/
/* Index: index_ppm_tst_case_group_app_version_id               */
/*==============================================================*/
create index index_ppm_tst_case_group_app_version_id on ppm_tst_case_group
(
   app_version_id
);

/*==============================================================*/
/* Table: ppm_tst_case_step                                     */
/*==============================================================*/
create table if not exists ppm_tst_case_step
(
   id                   bigint not null,
   org_id               bigint not null default 0,
   app_id               bigint not null default 0,
   app_version_id       bigint not null default 0,
   case_id              bigint not null default 0,
   step_desc            text,
   expect_desc          text,
   creator              bigint not null default 0,
   create_time          datetime not null default CURRENT_TIMESTAMP,
   updator              bigint not null default 0,
   update_time          datetime not null default CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
   version              int not null default 1,
   is_delete            tinyint not null default 2,
   primary key (id)
);

/*==============================================================*/
/* Index: index_ppm_tst_case_org_id                             */
/*==============================================================*/
create index index_ppm_tst_case_org_id on ppm_tst_case_step
(
   org_id
);

/*==============================================================*/
/* Index: index_ppm_tst_case_create_time                        */
/*==============================================================*/
create index index_ppm_tst_case_create_time on ppm_tst_case_step
(
   create_time
);

/*==============================================================*/
/* Index: index_ppm_tst_case_id                                 */
/*==============================================================*/
create index index_ppm_tst_case_id on ppm_tst_case_step
(
   case_id
);

/*==============================================================*/
/* Index: index_ppm_tst_pro_id                                  */
/*==============================================================*/
create index index_ppm_tst_pro_id on ppm_tst_case_step
(
   app_id
);

/*==============================================================*/
/* Index: index_ppm_tst_app_version_id                          */
/*==============================================================*/
create index index_ppm_tst_app_version_id on ppm_tst_case_step
(
   app_version_id
);

/*==============================================================*/
/* Table: ppm_tst_plan                                          */
/*==============================================================*/
create table if not exists ppm_tst_plan
(
   id                   bigint not null,
   org_id               bigint not null default 0,
   app_id               bigint not null default 0,
   app_version_id       bigint not null default 0,
   name                 varchar(128) not null default '',
   iter_id              bigint not null default 0,
   default_tester_id    bigint not null default 0,
   plan_desc            text,
   status               tinyint not null default 1,
   selective_type       tinyint(2) not null default 1,
   creator              bigint not null default 0,
   create_time          datetime not null default CURRENT_TIMESTAMP,
   updator              bigint not null default 0,
   update_time          datetime not null default CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
   version              int not null default 1,
   is_delete            tinyint not null default 2,
   primary key (id)
);

/*==============================================================*/
/* Index: index_ppm_tst_case_group_org_id                       */
/*==============================================================*/
create index index_ppm_tst_case_group_org_id on ppm_tst_plan
(
   org_id
);

/*==============================================================*/
/* Index: index_ppm_tst_case_group_create_time                  */
/*==============================================================*/
create index index_ppm_tst_case_group_create_time on ppm_tst_plan
(
   create_time
);

/*==============================================================*/
/* Index: index_ppm_tst_case_group_default_tester_id            */
/*==============================================================*/
create index index_ppm_tst_case_group_default_tester_id on ppm_tst_plan
(
   default_tester_id
);

/*==============================================================*/
/* Index: index_ppm_tst_case_group_app_id                       */
/*==============================================================*/
create index index_ppm_tst_case_group_app_id on ppm_tst_plan
(
   app_id
);

/*==============================================================*/
/* Index: index_ppm_tst_case_group_iter_id                      */
/*==============================================================*/
create index index_ppm_tst_case_group_iter_id on ppm_tst_plan
(
   iter_id
);

/*==============================================================*/
/* Index: index_ppm_tst_case_group_app_version_id               */
/*==============================================================*/
create index index_ppm_tst_case_group_app_version_id on ppm_tst_plan
(
   app_version_id
);

/*==============================================================*/
/* Table: ppm_tst_plan_case                                     */
/*==============================================================*/
create table if not exists ppm_tst_plan_case
(
   id                   bigint not null,
   org_id               bigint not null default 0,
   app_id               bigint not null default 0,
   app_version_id       bigint not null default 0,
   case_id              bigint not null default 0,
   status               tinyint not null default 1,
   tester_id            bigint not null default 0,
   creator              bigint not null default 0,
   create_time          datetime not null default CURRENT_TIMESTAMP,
   updator              bigint not null default 0,
   update_time          datetime not null default CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
   version              int not null default 1,
   is_delete            tinyint not null default 2,
   primary key (id)
);

/*==============================================================*/
/* Index: index_ppm_tst_case_org_id                             */
/*==============================================================*/
create index index_ppm_tst_case_org_id on ppm_tst_plan_case
(
   org_id
);

/*==============================================================*/
/* Index: index_ppm_tst_case_create_time                        */
/*==============================================================*/
create index index_ppm_tst_case_create_time on ppm_tst_plan_case
(
   create_time
);

/*==============================================================*/
/* Index: index_ppm_tst_case_id                                 */
/*==============================================================*/
create index index_ppm_tst_case_id on ppm_tst_plan_case
(
   case_id
);

/*==============================================================*/
/* Index: index_ppm_tst_pro_id                                  */
/*==============================================================*/
create index index_ppm_tst_pro_id on ppm_tst_plan_case
(
   app_id
);

/*==============================================================*/
/* Index: index_ppm_tst_app_version_id                          */
/*==============================================================*/
create index index_ppm_tst_app_version_id on ppm_tst_plan_case
(
   app_version_id
);

/*==============================================================*/
/* Table: ppm_tst_plan_case_issue                               */
/*==============================================================*/
create table if not exists ppm_tst_plan_case_issue
(
   id                   bigint not null,
   org_id               bigint not null default 0,
   app_id               bigint not null default 0,
   app_version_id       bigint not null default 0,
   plan_case_id         bigint not null default 0,
   issue_id             bigint not null default 0,
   case_step_id         bigint not null default 0,
   case_id              bigint not null,
   creator              bigint not null default 0,
   create_time          datetime not null default CURRENT_TIMESTAMP,
   updator              bigint not null default 0,
   update_time          datetime not null default CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
   version              int not null default 1,
   is_delete            tinyint not null default 2,
   primary key (id)
);

/*==============================================================*/
/* Index: index_ppm_tst_case_org_id                             */
/*==============================================================*/
create index index_ppm_tst_case_org_id on ppm_tst_plan_case_issue
(
   org_id
);

/*==============================================================*/
/* Index: index_ppm_tst_case_create_time                        */
/*==============================================================*/
create index index_ppm_tst_case_create_time on ppm_tst_plan_case_issue
(
   create_time
);

/*==============================================================*/
/* Index: index_ppm_tst_pro_id                                  */
/*==============================================================*/
create index index_ppm_tst_pro_id on ppm_tst_plan_case_issue
(
   app_id
);

/*==============================================================*/
/* Index: index_ppm_tst_plan_case_id                            */
/*==============================================================*/
create index index_ppm_tst_plan_case_id on ppm_tst_plan_case_issue
(
   plan_case_id
);

/*==============================================================*/
/* Index: index_ppm_tst_issue_id                                */
/*==============================================================*/
create index index_ppm_tst_issue_id on ppm_tst_plan_case_issue
(
   issue_id
);

/*==============================================================*/
/* Index: index_ppm_tst_case_id                                 */
/*==============================================================*/
create index index_ppm_tst_case_id on ppm_tst_plan_case_issue
(
   case_id
);

/*==============================================================*/
/* Index: index_ppm_tst_app_version_id                          */
/*==============================================================*/
create index index_ppm_tst_app_version_id on ppm_tst_plan_case_issue
(
   app_version_id
);

/*==============================================================*/
/* Table: ppm_tst_plan_case_step                                */
/*==============================================================*/
create table if not exists ppm_tst_plan_case_step
(
   id                   bigint not null,
   org_id               bigint not null default 0,
   app_id               bigint not null default 0,
   app_version_id       bigint not null default 0,
   plan_case_id         bigint not null default 0,
   case_id              bigint not null default 0,
   case_step_id         bigint not null default 0,
   status               tinyint not null default 1,
   actual_result_desc   text,
   creator              bigint not null default 0,
   create_time          datetime not null default CURRENT_TIMESTAMP,
   updator              bigint not null default 0,
   update_time          datetime not null default CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
   version              int not null default 1,
   is_delete            tinyint not null default 2,
   primary key (id)
);

/*==============================================================*/
/* Index: index_ppm_tst_case_org_id                             */
/*==============================================================*/
create index index_ppm_tst_case_org_id on ppm_tst_plan_case_step
(
   org_id
);

/*==============================================================*/
/* Index: index_ppm_tst_case_create_time                        */
/*==============================================================*/
create index index_ppm_tst_case_create_time on ppm_tst_plan_case_step
(
   create_time
);

/*==============================================================*/
/* Index: index_ppm_tst_case_id                                 */
/*==============================================================*/
create index index_ppm_tst_case_id on ppm_tst_plan_case_step
(
   case_id
);

/*==============================================================*/
/* Index: index_ppm_tst_plan_case_id                            */
/*==============================================================*/
create index index_ppm_tst_plan_case_id on ppm_tst_plan_case_step
(
   plan_case_id
);

/*==============================================================*/
/* Index: index_ppm_tst_pro_id                                  */
/*==============================================================*/
create index index_ppm_tst_pro_id on ppm_tst_plan_case_step
(
   app_id
);

/*==============================================================*/
/* Index: index_ppm_tst_app_version_id                          */
/*==============================================================*/
create index index_ppm_tst_app_version_id on ppm_tst_plan_case_step
(
   app_version_id
);
