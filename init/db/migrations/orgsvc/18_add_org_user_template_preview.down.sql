delete from `ppm_org_organization` where id=999;

delete from `ppm_org_user` where id in (2, 10);

delete from `ppm_org_user_out_info` where id in (999, 1000);

delete from `ppm_org_organization_out_info` where id=999;

delete from `ppm_org_user_organization` where id in (998, 1000);

delete from `ppm_orc_config` where id=1000;
