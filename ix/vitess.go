package ix

import (
	"encoding/json"
	"fmt"
	"github.com/blastrain/vitess-sqlparser/tidbparser/ast"
	"github.com/blastrain/vitess-sqlparser/tidbparser/parser"
	"testing"
)

func JsonDump(o any) string {
	b, err := json.MarshalIndent(o, "", "  ")
	if err != nil {
		return ""
	}
	return string(b)
}

type DropColumnVisitor struct {
	TableName string
	Columns   []string
}

func (v *DropColumnVisitor) IsValid() bool {
	return v.TableName != "" && len(v.Columns) > 0
}

func (v *DropColumnVisitor) Enter(in ast.Node) (out ast.Node, skipChildren bool) {
	fmt.Printf("%T\n", in)
	switch in.(type) {
	case *ast.AlterTableStmt:
		if v.TableName != "" {
			break
		}
		node := in.(*ast.AlterTableStmt)
		v.TableName = node.Table.Name.String()
		for _, spec := range node.Specs {
			if spec.Tp == ast.AlterTableDropColumn {
				v.Columns = append(v.Columns, spec.OldColumnName.OrigColName())
			}
		}
	default:
		break
	}
	return in, true
}

func (v *DropColumnVisitor) Leave(in ast.Node) (out ast.Node, ok bool) {
	return in, true
}

func TestParseDropColumn(t *testing.T) {
	sql := "ALTER TABLE Persons DROP COLUMN DateOfBirth, DROP COLUMN ID;"
	sqlParser := parser.New()
	stmtNodes, err := sqlParser.Parse(sql, "", "")
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}
	t.Logf("stmt: %s", JsonDump(stmtNodes))

	v := &DropColumnVisitor{
		TableName: "",
		Columns:   []string{},
	}
	for _, stmtNode := range stmtNodes {
		stmtNode.Accept(v)
	}
	t.Logf("visitor: %s", JsonDump(v))

	if !v.IsValid() {
		t.Fatalf("invalid drop column ddl")
	}
	t.Logf("drop columns %v at table %s", v.Columns, v.TableName)
}
