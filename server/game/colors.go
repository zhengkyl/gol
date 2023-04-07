package game

type playerColor struct {
	cursor string
	cell   string
}

var ColorTable = [11]playerColor{
	{
		"#080808",
		"0",
	},
	{
		"#ff0000",
		"#ff5f5f",
	},
	{
		"#d75f00",
		"#ff8700",
	},
	{
		"#ffd700",
		"#ffff5f",
	},
	{
		"#87af00",
		"#afff00",
	},
	{
		"#005f00",
		"#00d700",
	},
	{
		"#00afff",
		"#00ffff",
	},
	{
		"#005f87",
		"#0087ff",
	},
	{
		"#d700ff",
		"#d787ff",
	},
	{
		"#ff00af",
		"#ff5faf",
	},
	{
		"#afafd7",
		"#eeeeee",
	},
}
