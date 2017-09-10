package jexweb

Config struct {
	Address       string `flag:"|:8080|http listening Address"`
	AssetsDir     string `flag:"|public/assets|website Assets storage path"`
	PublicDir     string `flag:"|public|website Public storage path"`
	TemplateDir   string `flag:"|templates|website Templates storage path"`
	AppLayout     string `flag:"|layout|website global layout template"`
	LoginPath     string `flag:"|/login|website login url"`
	IsDevelopment bool   //调试模式
}

// func newConfig(config ...map[string]interface{}) *Config {
// 	cnf := &Config{
// 		Web: jexConfig{
// 			Address:       ":8080",
// 			PublicDir:     "public",
// 			AssetsDir:     "public/assets",
// 			TemplateDir:   "templates",
// 			AppLayout:     "layout",
// 			LoginPath:     "/login",
// 			IsDevelopment: true,
// 		},
// 	}
// 	if len(config) > 0 {
// 		cnf.ExtendConfig = config[0]
// 	}
// 	return cnf
// }

// func (cfg *Config) Save() error {
// 	return cfg.loadConfig(cfg._filename)
// }

// func (cfg *Config) loadConfig(filename ...string) error {
// 	if len(filename) == 0 {
// 		//创建默认配置文件cnf
// 		fullExeFilename, _ := exec.LookPath(os.Args[0])
// 		fullPath := filepath.Dir(fullExeFilename)
// 		filename = append(filename, filepath.Join(fullPath, "cnf"))
// 	}
// 	cfg._filename = filename[0]
// 	if finfo, err := os.Stat(filename[0]); err != nil || finfo.IsDir() {
// 		//配置文件不存在，根据对象中默认参数创建
// 		saveData, _ := json.Marshal(cfg)
// 		//json格式化
// 		var out bytes.Buffer
// 		json.Indent(&out, saveData, "", "\t")
// 		file, _ := os.Create(filename[0])
// 		_, err := out.WriteTo(file)
// 		//	err := ioutil.WriteFile(filename[0], &out, os.ModeAppend)
// 		if err != nil {
// 			return err
// 		}
// 		return nil
// 	}
// 	if bts, err := ioutil.ReadFile(filename[0]); err == nil {
// 		_cfg := &Config{}
// 		json.Unmarshal(bts, _cfg)
// 		cfg.Web = _cfg.Web
// 		for k, v := range cfg.ExtendConfig {
// 			cfg.getExtendConfig(_cfg, k, v)
// 		}
// 		return nil
// 	} else {
// 		return err
// 	}
// }

// func (cfg *Config) getExtendConfig(_cfg *Config, key string, v interface{}) {
// 	bytes, _ := json.Marshal(_cfg.ExtendConfig[key])
// 	json.Unmarshal(bytes, v)
// }
