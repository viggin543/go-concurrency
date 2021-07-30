package lowlevel

func DoWork() (string,error){
	return "", &Error{}
}

type Error struct { error }


