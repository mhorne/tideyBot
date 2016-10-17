package soundPlayer


// Array of all the sounds we have
var AIRHORN *soundCollection = &soundCollection{
	prefix: "airhorn",
	commands: []string{
		"!airhorn",
	},
	sounds: []*sound{
		createSound("default", 1000, 250),
		createSound("reverb", 800, 250),
		createSound("spam", 800, 0),
		createSound("tripletap", 800, 250),
		createSound("fourtap", 800, 250),
		createSound("distant", 500, 250),
		createSound("echo", 500, 250),
		createSound("clownfull", 250, 250),
		createSound("clownshort", 250, 250),
		createSound("clownspam", 250, 0),
		createSound("highfartlong", 200, 250),
		createSound("highfartshort", 200, 250),
		createSound("midshort", 100, 250),
		createSound("truck", 10, 250),
	},
}

var KHALED *soundCollection = &soundCollection{
	prefix:    "another",
	chainWith: AIRHORN,
	commands: []string{
		"!anotha",
		"!anothaone",
	},
	sounds: []*sound{
		createSound("one", 1, 250),
		createSound("one_classic", 1, 250),
		createSound("one_echo", 1, 250),
	},
}

var CENA *soundCollection = &soundCollection{
	prefix: "jc",
	commands: []string{
		"!johncena",
		"!cena",
	},
	sounds: []*sound{
		createSound("airhorn", 1, 250),
		createSound("echo", 1, 250),
		createSound("full", 1, 250),
		createSound("jc", 1, 250),
		createSound("nameis", 1, 250),
		createSound("spam", 1, 250),
	},
}

var ETHAN *soundCollection = &soundCollection{
	prefix: "ethan",
	commands: []string{
		"!ethan",
		"!eb",
		"!ethanbradberry",
		"!h3h3",
	},
	sounds: []*sound{
		createSound("areyou_classic", 100, 250),
		createSound("areyou_condensed", 100, 250),
		createSound("areyou_crazy", 100, 250),
		createSound("areyou_ethan", 100, 250),
		createSound("classic", 100, 250),
		createSound("echo", 100, 250),
		createSound("high", 100, 250),
		createSound("slowandlow", 100, 250),
		createSound("cuts", 30, 250),
		createSound("beat", 30, 250),
		createSound("sodiepop", 1, 250),
	},
}

var COW *soundCollection = &soundCollection{
	prefix: "cow",
	commands: []string{
		"!stan",
		"!stanislav",
	},
	sounds: []*sound{
		createSound("herd", 10, 250),
		createSound("moo", 10, 250),
		createSound("x3", 1, 250),
	},
}

var BIRTHDAY *soundCollection = &soundCollection{
	prefix: "birthday",
	commands: []string{
		"!birthday",
		"!bday",
	},
	sounds: []*sound{
		createSound("horn", 50, 250),
		createSound("horn3", 30, 250),
		createSound("sadhorn", 25, 250),
		createSound("weakhorn", 25, 250),
	},
}

var WOW *soundCollection = &soundCollection{
	prefix: "wow",
	commands: []string{
		"!wowthatscool",
		"!wtc",
	},
	sounds: []*sound{
		createSound("thatscool", 50, 250),
	},
}

var HYPE *soundCollection = &soundCollection{
	prefix: "hype",
	commands: []string{
		"!hype",
		"!na",
		"!cs",
	},
	sounds: []*sound{
		createSound("bestteam", 100, 250),
		createSound("brabrabravo", 100, 250),
		createSound("brabrabravo2", 100, 250),
		createSound("cans", 100, 250),
		createSound("cans2", 100, 250),
		createSound("givenoise", 100, 250),
		createSound("givenoise2", 100, 250),
		createSound("givenoise3", 100, 250),
		createSound("goingtowar", 100, 250),
		createSound("herewegoagain", 100, 250),
		createSound("millions", 100, 250),
		createSound("poland", 100, 250),
		createSound("ready", 100, 250),
		createSound("show", 100, 250),
		createSound("veryexciting", 100, 250),
		createSound("warisover", 100, 250),
	},
}

var NUTSHACK *soundCollection = &soundCollection{
	prefix: "nut",
	commands: []string{
		"!nutshack",
	},
	sounds: []*sound{
		createSound("shack", 100, 250),
		createSound("dick", 100, 250),
	},
}

var GOLF *soundCollection = &soundCollection{
	prefix: "golf",
	commands: []string{
		"!golf",
		"!gwf",
	},
	sounds: []*sound{
		createSound("gwf", 100, 250),
	},
}

var IOD44 *soundCollection = &soundCollection{
	prefix: "de",
	commands: []string{
		"!44",
		"!de44",
		"!IOD44",
	},
	sounds: []*sound{
		createSound("44", 100, 250),
	},
}

var RICKANDMORTY *soundCollection = &soundCollection{
	prefix: "rick",
	commands: []string{
		"!myman",
		"!yes",
	},
	sounds: []*sound{
		createSound("myman", 100, 250),
		createSound("yes", 100, 250),
	},
}

var MEME *soundCollection = &soundCollection{
	prefix: "meme",
	commands: []string{
		"!meme",
	},
	sounds: []*sound{
		createSound("prettygood", 100, 250),
		createSound("gay", 100, 250),
		createSound("idunnoshort", 100, 250),
		createSound("idunnowine", 100, 250),
	},
}

var COLLECTIONS []*soundCollection = []*soundCollection{
	AIRHORN,
	KHALED,
	CENA,
	ETHAN,
	COW,
	BIRTHDAY,
	WOW,
	HYPE,
	NUTSHACK,
	GOLF,
	IOD44,
	RICKANDMORTY,
	MEME,
}
