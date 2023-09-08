package main

type ShortListItem struct {
	contact		Contact
	queried 	bool
	responded	bool
}

type ShortList struct {
	list 	[]ShortListItem
	//TODO mutex
}

func (shortList *ShortList) Fill(contact *[]Contact) {

}