update ppm_prs_process_status set bg_style = "#DBDBDB" where lang_code = "ProcessStatus.Iteration.NotStart";
update ppm_prs_process_status set bg_style = "#FFCD1C" where lang_code = "ProcessStatus.Iteration.Running";
update ppm_prs_process_status set font_style = "#69A922" where lang_code = "ProcessStatus.Iteration.Complete";
