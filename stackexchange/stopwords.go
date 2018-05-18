package stackexchange

// Words we consider not to be useful as tech tags, i.e., in prose they are probably intended as plain English.
// Judgement call! We are betting that lemmatizing these words is undesirable or useless.
// Inspried by https://github.com/first20hours/google-10000-english
var stopWords = map[string]struct{}{
	"the":           exists,
	"of":            exists,
	"and":           exists,
	"to":            exists,
	"a":             exists,
	"in":            exists,
	"for":           exists,
	"is":            exists,
	"on":            exists,
	"that":          exists,
	"by":            exists,
	"this":          exists,
	"with":          exists,
	"i":             exists,
	"you":           exists,
	"it":            exists,
	"not":           exists,
	"or":            exists,
	"be":            exists,
	"are":           exists,
	"from":          exists,
	"at":            exists,
	"as":            exists,
	"your":          exists,
	"all":           exists,
	"have":          exists,
	"new":           exists,
	"more":          exists,
	"an":            exists,
	"was":           exists,
	"we":            exists,
	"will":          exists,
	"home":          exists,
	"can":           exists,
	"us":            exists,
	"about":         exists,
	"if":            exists,
	"page":          exists,
	"my":            exists,
	"has":           exists,
	"search":        exists,
	"free":          exists,
	"but":           exists,
	"our":           exists,
	"one":           exists,
	"other":         exists,
	"do":            exists,
	"no":            exists,
	"information":   exists,
	"time":          exists,
	"they":          exists,
	"site":          exists,
	"he":            exists,
	"up":            exists,
	"may":           exists,
	"what":          exists,
	"which":         exists,
	"their":         exists,
	"news":          exists,
	"out":           exists,
	"use":           exists,
	"any":           exists,
	"there":         exists,
	"see":           exists,
	"only":          exists,
	"so":            exists,
	"his":           exists,
	"when":          exists,
	"contact":       exists,
	"here":          exists,
	"business":      exists,
	"who":           exists,
	"web":           exists,
	"also":          exists,
	"now":           exists,
	"help":          exists,
	"get":           exists,
	"pm":            exists,
	"view":          exists,
	"online":        exists,
	"first":         exists,
	"am":            exists,
	"been":          exists,
	"would":         exists,
	"how":           exists,
	"were":          exists,
	"me":            exists,
	"services":      exists,
	"some":          exists,
	"these":         exists,
	"click":         exists,
	"its":           exists,
	"like":          exists,
	"service":       exists,
	"than":          exists,
	"find":          exists,
	"price":         exists,
	"date":          exists,
	"back":          exists,
	"top":           exists,
	"people":        exists,
	"had":           exists,
	"list":          exists,
	"name":          exists,
	"just":          exists,
	"over":          exists,
	"state":         exists,
	"year":          exists,
	"day":           exists,
	"into":          exists,
	"email":         exists,
	"two":           exists,
	"health":        exists,
	"world":         exists,
	"re":            exists,
	"next":          exists,
	"used":          exists,
	"work":          exists,
	"last":          exists,
	"most":          exists,
	"products":      exists,
	"music":         exists,
	"buy":           exists,
	"data":          exists,
	"make":          exists,
	"them":          exists,
	"should":        exists,
	"product":       exists,
	"system":        exists,
	"post":          exists,
	"her":           exists,
	"city":          exists,
	"add":           exists,
	"number":        exists,
	"such":          exists,
	"please":        exists,
	"available":     exists,
	"copyright":     exists,
	"support":       exists,
	"message":       exists,
	"after":         exists,
	"best":          exists,
	"software":      exists,
	"then":          exists,
	"good":          exists,
	"video":         exists,
	"well":          exists,
	"where":         exists,
	"info":          exists,
	"public":        exists,
	"books":         exists,
	"through":       exists,
	"each":          exists,
	"she":           exists,
	"review":        exists,
	"years":         exists,
	"order":         exists,
	"very":          exists,
	"book":          exists,
	"items":         exists,
	"company":       exists,
	"read":          exists,
	"group":         exists,
	"sex":           exists,
	"need":          exists,
	"many":          exists,
	"user":          exists,
	"said":          exists,
	"does":          exists,
	"set":           exists,
	"under":         exists,
	"general":       exists,
	"mail":          exists,
	"full":          exists,
	"map":           exists,
	"reviews":       exists,
	"program":       exists,
	"life":          exists,
	"know":          exists,
	"games":         exists,
	"way":           exists,
	"days":          exists,
	"management":    exists,
	"part":          exists,
	"could":         exists,
	"great":         exists,
	"united":        exists,
	"hotel":         exists,
	"real":          exists,
	"item":          exists,
	"international": exists,
	"center":        exists,
	"must":          exists,
	"store":         exists,
	"travel":        exists,
	"comments":      exists,
	"made":          exists,
	"development":   exists,
	"report":        exists,
	"off":           exists,
	"member":        exists,
	"details":       exists,
	"line":          exists,
	"terms":         exists,
	"before":        exists,
	"did":           exists,
	"send":          exists,
	"right":         exists,
	"type":          exists,
	"because":       exists,
	"local":         exists,
	"those":         exists,
	"using":         exists,
	"results":       exists,
	"office":        exists,
	"education":     exists,
	"national":      exists,
	"car":           exists,
	"design":        exists,
	"take":          exists,
	"posted":        exists,
	"internet":      exists,
	"address":       exists,
	"community":     exists,
	"within":        exists,
	"states":        exists,
	"area":          exists,
	"want":          exists,
	"phone":         exists,
	"dvd":           exists,
	"shipping":      exists,
	"reserved":      exists,
	"subject":       exists,
	"between":       exists,
	"forum":         exists,
	"family":        exists,
	"long":          exists,
	"based":         exists,
	"code":          exists,
	"show":          exists,
	"even":          exists,
	"black":         exists,
	"check":         exists,
	"special":       exists,
	"prices":        exists,
	"website":       exists,
	"index":         exists,
	"being":         exists,
	"women":         exists,
	"much":          exists,
	"sign":          exists,
	"file":          exists,
	"link":          exists,
	"open":          exists,
	"today":         exists,
	"technology":    exists,
	"south":         exists,
	"case":          exists,
	"project":       exists,
	"same":          exists,
	"pages":         exists,
	"version":       exists,
	"section":       exists,
	"own":           exists,
	"found":         exists,
	"sports":        exists,
	"house":         exists,
	"related":       exists,
	"security":      exists,
	"both":          exists,
	"county":        exists,
	"american":      exists,
	"photo":         exists,
	"game":          exists,
	"members":       exists,
	"power":         exists,
	"while":         exists,
	"care":          exists,
	"network":       exists,
	"down":          exists,
	"computer":      exists,
	"systems":       exists,
	"three":         exists,
	"total":         exists,
	"place":         exists,
	"end":           exists,
	"following":     exists,
	"download":      exists,

	"him":         exists,
	"without":     exists,
	"per":         exists,
	"access":      exists,
	"think":       exists,
	"north":       exists,
	"resources":   exists,
	"current":     exists,
	"posts":       exists,
	"big":         exists,
	"media":       exists,
	"law":         exists,
	"control":     exists,
	"water":       exists,
	"history":     exists,
	"pictures":    exists,
	"size":        exists,
	"art":         exists,
	"personal":    exists,
	"since":       exists,
	"including":   exists,
	"guide":       exists,
	"shop":        exists,
	"board":       exists,
	"location":    exists,
	"change":      exists,
	"white":       exists,
	"text":        exists,
	"small":       exists,
	"rating":      exists,
	"rate":        exists,
	"government":  exists,
	"children":    exists,
	"during":      exists,
	"usa":         exists,
	"return":      exists,
	"students":    exists,
	"shopping":    exists,
	"account":     exists,
	"times":       exists,
	"sites":       exists,
	"level":       exists,
	"digital":     exists,
	"profile":     exists,
	"previous":    exists,
	"form":        exists,
	"events":      exists,
	"love":        exists,
	"old":         exists,
	"john":        exists,
	"main":        exists,
	"call":        exists,
	"hours":       exists,
	"image":       exists,
	"department":  exists,
	"title":       exists,
	"description": exists,
	"non":         exists,
	"insurance":   exists,
	"another":     exists,
	"why":         exists,
	"shall":       exists,
	"property":    exists,
	"class":       exists,
	"cd":          exists,
	"still":       exists,
	"money":       exists,
	"quality":     exists,
	"every":       exists,
	"listing":     exists,
	"content":     exists,
	"country":     exists,
	"private":     exists,
	"little":      exists,
	"visit":       exists,
	"save":        exists,
	"tools":       exists,
	"low":         exists,
	"reply":       exists,
	"customer":    exists,
	"december":    exists,
	"compare":     exists,
	"movies":      exists,
	"include":     exists,
	"college":     exists,
	"value":       exists,
	"article":     exists,
	"york":        exists,
	"man":         exists,
	"card":        exists,
	"jobs":        exists,
	"provide":     exists,
	"food":        exists,
	"source":      exists,
	"author":      exists,
	"different":   exists,
	"press":       exists,
	"learn":       exists,
	"sale":        exists,
	"around":      exists,
	"print":       exists,
	"course":      exists,
	"job":         exists,
	"canada":      exists,
	"process":     exists,
	"teen":        exists,
	"room":        exists,
	"stock":       exists,
	"training":    exists,
	"too":         exists,
	"credit":      exists,
	"point":       exists,
	"join":        exists,
	"science":     exists,
	"men":         exists,
	"categories":  exists,
	"advanced":    exists,
	"west":        exists,
	"sales":       exists,
	"look":        exists,
	"english":     exists,
	"left":        exists,
	"team":        exists,
	"estate":      exists,
	"box":         exists,
	"conditions":  exists,
	"select":      exists,
	"windows":     exists,
	"photos":      exists,
	"gay":         exists,
	"thread":      exists,
	"week":        exists,
	"category":    exists,
	"note":        exists,
	"live":        exists,
	"large":       exists,
	"gallery":     exists,
	"table":       exists,
	"register":    exists,
	"however":     exists,
	"june":        exists,
	"october":     exists,
	"november":    exists,
	"market":      exists,
	"library":     exists,
	"really":      exists,
	"action":      exists,
	"start":       exists,
	"series":      exists,
	"model":       exists,
	"features":    exists,
	"air":         exists,
	"industry":    exists,
	"plan":        exists,
	"human":       exists,
	"provided":    exists,
	"tv":          exists,
	"yes":         exists,
	"required":    exists,
	"second":      exists,
	"hot":         exists,
	"accessories": exists,
	"cost":        exists,
	"movie":       exists,
	"forums":      exists,
	"march":       exists,
	"september":   exists,
	"better":      exists,
	"say":         exists,
	"questions":   exists,
	"july":        exists,
	"yahoo":       exists,
	"going":       exists,
	"medical":     exists,
	"test":        exists,
	"friend":      exists,
	"come":        exists,
	"dec":         exists,
	"server":      exists,
	"function":    exists,
}

func isStopWord(s string) bool {
	_, exists := stopWords[s]
	return exists
}
