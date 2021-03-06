package anaconda

import (
	"net/url"
	"strconv"
	"strings"
)

// CreateList implements /lists/create.json
func (a TwitterApi) CreateList(name, description string, v url.Values) (list List, err error) {
	v = cleanValues(v)
	v.Set("name", name)
	v.Set("description", description)

	response_ch := make(chan response)
	a.queryQueue <- query{a.baseUrl + "/lists/create.json", v, &list, _POST, response_ch}
	return list, (<-response_ch).err
}

// AddUserToList implements /lists/members/create.json
func (a TwitterApi) AddUserToList(screenName string, listID int64, v url.Values) (users []User, err error) {
	v = cleanValues(v)
	v.Set("list_id", strconv.FormatInt(listID, 10))
	v.Set("screen_name", screenName)

	var addUserToListResponse AddUserToListResponse

	response_ch := make(chan response)
	a.queryQueue <- query{a.baseUrl + "/lists/members/create.json", v, &addUserToListResponse, _POST, response_ch}
	return addUserToListResponse.Users, (<-response_ch).err
}

// AddMultipleUsersToList implements /lists/members/create_all.json
func (a TwitterApi) AddMultipleUsersToList(screenNames []string, listID int64, v url.Values) (list List, err error) {
	v = cleanValues(v)
	v.Set("list_id", strconv.FormatInt(listID, 10))
	v.Set("screen_name", strings.Join(screenNames, ","))

	response_ch := make(chan response)
	a.queryQueue <- query{a.baseUrl + "/lists/members/create_all.json", v, &list, _POST, response_ch}
	r := <-response_ch
	return list, r.err
}

func (a TwitterApi) GetListMembers(listID int64, v url.Values) ([]User, error) {
	v = cleanValues(v)
	v.Set("list_id", strconv.FormatInt(listID, 10))
	v.Set("cursor", "-1")
 
	var all []User
	var lmResp ListMembershipResponse
	response_ch := make(chan response)
	a.queryQueue <- query{a.baseUrl + "/lists/members.json", v, &lmResp, _GET, response_ch}
	if err := (<-response_ch).err; err != nil {
		return nil, err
	}
	all = lmResp.Users
	for lmResp.NextCursor != 0 {
		//log.Printf("additonal list(%d) members req. %d %d.", listID, lmResp.NextCursor, len(all))
 
		v.Set("cursor", strconv.FormatInt(int64(lmResp.NextCursor), 10))
 
		response_ch = make(chan response)
		a.queryQueue <- query{a.baseUrl + "/lists/members.json", v, &lmResp, _GET, response_ch}
		if err := (<-response_ch).err; err != nil {
			return nil, err
		}
		all = append(all, lmResp.Users...)
	}
	return all, nil
}

// GetListsOwnedBy implements /lists/ownerships.json
// screen_name, count, and cursor are all optional values
func (a TwitterApi) GetListsOwnedBy(userID int64, v url.Values) (lists []List, err error) {
	v = cleanValues(v)
	v.Set("user_id", strconv.FormatInt(userID, 10))

	var listResponse ListResponse

	response_ch := make(chan response)
	a.queryQueue <- query{a.baseUrl + "/lists/ownerships.json", v, &listResponse, _GET, response_ch}
	return listResponse.Lists, (<-response_ch).err
}

func (a TwitterApi) GetListTweets(listID int64, includeRTs bool, v url.Values) (tweets []Tweet, err error) {
	v = cleanValues(v)
	v.Set("list_id", strconv.FormatInt(listID, 10))
	v.Set("include_rts", strconv.FormatBool(includeRTs))

	response_ch := make(chan response)
	a.queryQueue <- query{a.baseUrl + "/lists/statuses.json", v, &tweets, _GET, response_ch}
	return tweets, (<-response_ch).err
}

// GetList implements /lists/show.json
func (a TwitterApi) GetList(listID int64, v url.Values) (list List, err error) {
	v = cleanValues(v)
	v.Set("list_id", strconv.FormatInt(listID, 10))

	response_ch := make(chan response)
	a.queryQueue <- query{a.baseUrl + "/lists/show.json", v, &list, _GET, response_ch}
	return list, (<-response_ch).err
}


// GetListBySlug implements /lists/show.json
func (a TwitterApi) GetListBySlug(slug string, ownerScreenName string, v url.Values) (list List, err error) {
	v = cleanValues(v)
	v.Set("slug", slug)
	v.Set("owner_screen_name", ownerScreenName)

	response_ch := make(chan response)
	a.queryQueue <- query{a.baseUrl + "/lists/show.json", v, &list, _GET, response_ch}
	return list, (<-response_ch).err
}

func (a TwitterApi) GetListTweetsBySlug(slug string, ownerScreenName string, includeRTs bool, v url.Values) (tweets []Tweet, err error) {
	v = cleanValues(v)
	v.Set("slug", slug)
	v.Set("owner_screen_name", ownerScreenName)
	v.Set("include_rts", strconv.FormatBool(includeRTs))

	response_ch := make(chan response)
	a.queryQueue <- query{a.baseUrl + "/lists/statuses.json", v, &tweets, _GET, response_ch}
	return tweets, (<-response_ch).err
}
