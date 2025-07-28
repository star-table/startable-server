drop table if exists ppm_prs_project_type_category;

delete from ppm_prs_project_type where id between 3 and 46;

delete from ppm_prs_project_object_type where id between 7 and 307;

delete from ppm_prs_project_type_project_object_type where id between 10 and 354;

ALTER TABLE `ppm_prs_project_type` DROP COLUMN `cover`;
