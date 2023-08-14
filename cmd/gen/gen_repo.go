package gen

import (
	"github.com/spf13/viper"
	"strings"

	"github.com/gogf/gf-cli/v2/library/mlog"
	"github.com/gogf/gf-cli/v2/library/utils"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gfile"
	"github.com/gogf/gf/v2/text/gstr"
)

func GenRepo(req GenReq) {
	context := gstr.ReplaceByMap(TempRepo, g.MapStrStr{
		"TempImportPkg":         viper.Get("server.go_module").(string),
		"TempSvcNameCaseCamel":  GetJsonTagFromCase(req.TableName, "Camel"),      // 表名 大驼峰
		"TempSvcNameCamelLower": GetJsonTagFromCase(req.TableName, "CamelLower"), // 表名 小驼峰
	})

	// 无符号转有符号
	var idSymbolConv string
	for _, v := range req.TableColumns {
		if v.Field == req.TableKey {
			idSymbolConv = generateStructFieldTypeNameForEntity(v)
		}
	}

	context = strings.Replace(context, "int64", idSymbolConv, -1)

	context = strings.Replace(context, "total uint64", "total int64", -1)

	path := req.RepositoryImplDir + "/repo_" + req.TableName + ".go"
	if err := gfile.PutContents(path, strings.TrimSpace(context)); err != nil {
		mlog.Fatalf("writing content to '%s' failed: %v", path, err)
	} else {
		utils.GoFmt(path)
		mlog.Print("generated:", path)
	}
}
