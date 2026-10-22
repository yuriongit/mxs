/*
Package cnf provides all the functionality
for MXS' config. Manages the application's
configuration and script directory.

Copyright © 2026 Yuri Okeren <yuri.dev44@outlook.com>.
*/
package cnf

type ScriptHeader = string

type ScriptHeaders struct {
	Description ScriptHeader
	Done        ScriptHeader
}
