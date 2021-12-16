package git

import (
	"strings"
	"testing"
)

func TestParseRepo(t *testing.T) {
	_, err := readCommits(strings.NewReader(commitsStr))
	if err != nil {
		t.Fatal(err)
	}

	//	for _, c := range cs {
	//		fmt.Println("%v", c)
	//	}
}

func TestParseRepo2(t *testing.T) {
	_, err := readCommits(strings.NewReader(commitStr2))
	if err != nil {
		t.Fatal(err)
	}

	//	for _, c := range cs {
	//		fmt.Println("%v", c)
	//	}
}

func TestParseCommit(t *testing.T) {
	_, err := readCommit(strings.NewReader(commitStr))
	if err != nil {
		t.Fatal(err)
	}

	//	for i, f := range c.Files {
	//		fmt.Println(i, f.Path)
	//		for k, d := range f.Diffs {
	//			fmt.Println(k, d.Unified)
	//			for j, l := range d.Lines {
	//				fmt.Println(j, l.Type, l.Text)
	//			}
	//		}
	//	}
}

var commitStr2 = `commit ffbe1443a318312a1a2917e71ef5c4a71c911a44 (tag: 7.0)
Author: Santi <santiago@syltek.com>
Date:   2015-03-27 00:05:52 +0100

    [CodeReview] modificar vista commits

diff --git a/Modules/CodeReview/CodeReview/Controllers/CommitsController.cs b/Modules/CodeReview/CodeReview/Controllers/CommitsController.cs
index d29e7f7..bc29635 100644
--- a/Modules/CodeReview/CodeReview/Controllers/CommitsController.cs
+++ b/Modules/CodeReview/CodeReview/Controllers/CommitsController.cs
@@ -21,7 +21,9 @@ namespace CodeReview
                                repoName = repos.Keys.First();
                        }
 
-                       var branch = this.Request["branch"] ?? "master";
+                       // mostrar la rama dev por defecto
+                       var branch = this.Request["branch"] ?? "dev";
+
                        model["branch"] = branch;
                        model["repoName"] = repoName;
                        model["unreadcomments"] = RepoUtil.UnreadCount(this.User);
@@ -46,11 +48,6 @@ namespace CodeReview
                        var code = this.GetString("searchCode");
                        var path = this.GetString("searchPath");
 
-                       if(this.GetValue<bool>("allbranches", false))
-                       {
-                               branch = "--all";
-                       }
-
                        var repository = new Repository(repo as string);
 
                        var commits = repository.Commits(start, pageSize, branch, commitMessage, code, path);
\ No newline at end of file
diff --git a/Modules/CodeReview/WebServer/Program.cs b/Modules/CodeReview/WebServer/Program.cs
index 136773a..f8f5c6d 100644
--- a/Modules/CodeReview/WebServer/Program.cs
+++ b/Modules/CodeReview/WebServer/Program.cs
@@ -28,8 +28,8 @@ namespace WebServer
                        var defaultParser = new SclRouteParser();
                        RouteParser.Default = defaultParser;
 
-            defaultParser.AddRoute("/", "codereview", "branches", "index");
-            defaultParser.AddRoute("/admin", "codereview", "branches", "index");
+                       defaultParser.AddRoute("/", "codereview", "commits", "index");
+                       defaultParser.AddRoute("/admin", "codereview", "commits", "index");
                        defaultParser.AddRoute("/codereview/commits/:repo/:sha1", "codereview", "commits", "commit");
                        defaultParser.AddRoute("/codereview/{controller}", "codereview", "{controller}", "index");
                        defaultParser.AddRoute("/codereview/{controller}/{action}/:", "codereview", "{controller}", "{action}");`

var commitStr = `commit ffbe1443a318312a1a2917e71ef5c4a71c911a44 (tag: 7.0)
Author: Santi <santiago@syltek.com>
Date:   2015-03-27 00:05:52 +0100

    [CodeReview] modificar vista commits

diff --git a/Modules/CodeReview/CodeReview/Controllers/CommitsController.cs b/Modules/CodeReview/CodeReview/Controllers/CommitsController.cs
index d29e7f7..bc29635 100644
--- a/Modules/CodeReview/CodeReview/Controllers/CommitsController.cs
+++ b/Modules/CodeReview/CodeReview/Controllers/CommitsController.cs
@@ -21,7 +21,9 @@ namespace CodeReview
                                repoName = repos.Keys.First();
                        }
 
-                       var branch = this.Request["branch"] ?? "master";
+                       // mostrar la rama dev por defecto
+                       var branch = this.Request["branch"] ?? "dev";
+
                        model["branch"] = branch;
                        model["repoName"] = repoName;
                        model["unreadcomments"] = RepoUtil.UnreadCount(this.User);
@@ -46,11 +48,6 @@ namespace CodeReview
                        var code = this.GetString("searchCode");
                        var path = this.GetString("searchPath");
 
-                       if(this.GetValue<bool>("allbranches", false))
-                       {
-                               branch = "--all";
-                       }
-
                        var repository = new Repository(repo as string);
 
                        var commits = repository.Commits(start, pageSize, branch, commitMessage, code, path);
diff --git a/Modules/CodeReview/WebServer/Program.cs b/Modules/CodeReview/WebServer/Program.cs
index 136773a..f8f5c6d 100644
--- a/Modules/CodeReview/WebServer/Program.cs
+++ b/Modules/CodeReview/WebServer/Program.cs
@@ -28,8 +28,8 @@ namespace WebServer
                        var defaultParser = new SclRouteParser();
                        RouteParser.Default = defaultParser;
 
-            defaultParser.AddRoute("/", "codereview", "branches", "index");
-            defaultParser.AddRoute("/admin", "codereview", "branches", "index");
+                       defaultParser.AddRoute("/", "codereview", "commits", "index");
+                       defaultParser.AddRoute("/admin", "codereview", "commits", "index");
                        defaultParser.AddRoute("/codereview/commits/:repo/:sha1", "codereview", "commits", "commit");
                        defaultParser.AddRoute("/codereview/{controller}", "codereview", "{controller}", "index");
                        defaultParser.AddRoute("/codereview/{controller}/{action}/:", "codereview", "{controller}", "{action}");`

var commitsStr = `commit 5f9b3a84a1f096d531f270e85aa66f4db8b8e449 (HEAD, tag: 8.0, dev)
Author: Santi <santiago@syltek.com>
Date:   2015-03-27 19:07:10 +0100

    [CustomerZone] Anular todas las reservas asociadas a una venta.
    
    Si no se cancelaran un usuario avispado podría pagar varias reservas,
    cancelar una, con lo que se le devuelve todo el dinero y
    quedarse con las otras gratis.

commit 54115c3a87992953986cd83fd5aebe1a6f5b2a78
Author: Santi <santiago@syltek.com>
Date:   2015-03-27 00:23:57 +0100

    [CodeReview] Mejorar estilos del nuevo listado de codereview

commit ffbe1443a318312a1a2917e71ef5c4a71c911a44 (HEAD, tag: test,23, tag: 7.0, master)
Author: Santi <santiago@syltek.com>
Date:   2015-03-27 00:05:52 +0100

    [CodeReview] modificar vista commits

commit a0249fa36b20f1227984ae143b50256ba8e76d94
Author: Santi <santiago@syltek.com>
Date:   2015-03-26 23:21:50 +0100

    Actualizar gitworkflow

commit eddbd46f462ac46b8f67a8353ba94a048fe36900 (syltek/dev)
Author: Santi <santiago@syltek.com>
Date:   2015-03-26 20:01:32 +0100

    añadir git workflow

commit dfc0fc38315c2a8aa5a8129167bd0324ee74b106 (padelmaster/master, jose/master)
Author: JNavero <jnavero@gmail.com>
Date:   2015-03-26 12:48:23 +0100

    Corregir Test para UserLogin (ExternalAuthenticators).
`
