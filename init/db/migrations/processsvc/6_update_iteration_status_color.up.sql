update ppm_prs_process_status set bg_style = "#377AFF" where lang_code = "ProcessStatus.Iteration.NotStart";
update ppm_prs_process_status set bg_style = "#F0A100" where lang_code = "ProcessStatus.Iteration.Running";
update ppm_prs_process_status set bg_style = "#25B47E" where lang_code = "ProcessStatus.Iteration.Complete";
