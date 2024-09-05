package utils

type Comment_Map struct {
	supports_multi bool
	single_comment string
	multi_comment_open string
	multi_comment_close string
}

/*
language comments Map, Key: extension name ,
value:Comment_Map struct
*/

var lookup_map map[string]Comment_Map = map[string]Comment_Map{
	"go":     {	true,"//","/*","*/"},
	"c":      {	true,"//","/*","*/"},               
	"cpp":    {	true,"//","/*","*/"},
	"js":     {	true,"//","/*","*/"},
	"ts":     {	true,"//","/*","*/"},
	"jsx":    {	true,"//","/*","*/"},
	"tsx":    {	true,"//","/*","*/"},
	"java":   {	true,"//","/*","*/"},
	"rs":     {	true,"//","/*","*/"},
	"swift":  {	true,"//","/*","*/"},
	"kt":     {	true,"//","/*","*/"},
	"php":    {	true,"//","/*","*/"},
	"m":      {	true,"//","/*","*/"},
	"groovy": {	true,"//","/*","*/"},
	"cs":     {	true,"//","/*","*/"},
	"scala":  {	true,"//","/*","*/"},
	"zig":    {	true,"//","/*","*/"},
	"gleam":  {	true,"//","/*","*/"},
	"py":     {	true,"#",`"""`,`"""`},
	"r":      {	false,"#","",""},                     
	"rb":     {	false,"#","",""},                   
	"sh":     {	true,"#",": '","'"},                   
	"pl":     {	true,"#","/*","*/"},                   
	"ex":     {	false,"#","",""},                   
	"exs":    {	false,"#","",""},                   
	"jl":     {	true,"#","#=","=#"},                   
	"lua":    {	true,"--","--[[","]]--"},  
	"hs":     {	true,"--","/*","*/"},                    
	"sql":    {	true,"--","/*","*/"},                    
	"cbl":    { true,"*","/*","*/"},                   
	"erl":    { true,"%","=begin","=cut"},                   
	"clj":    { false,";;","",""},                    
	"lisp":   { false,";","",""},                   
}