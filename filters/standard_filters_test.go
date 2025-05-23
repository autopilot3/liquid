package filters

import (
	"fmt"
	"os"
	"testing"
	"time"

	yaml "gopkg.in/yaml.v2"

	"github.com/stretchr/testify/require"

	"github.com/autopilot3/ap3-types-go/types/date"
	"github.com/autopilot3/liquid/expressions"
)

var filterTests = []struct {
	in       string
	expected interface{}
}{
	// value filters
	{`undefined | default: 2.99`, 2.99},
	{`nil | default: 2.99`, 2.99},
	{`false | default: 2.99`, 2.99},
	{`"" | default: 2.99`, 2.99},
	{`empty_array | default: 2.99`, 2.99},
	{`empty_map | default: 2.99`, 2.99},
	{`empty_map_slice | default: 2.99`, 2.99},
	{`true | default: 2.99`, true},
	{`"true" | default: 2.99`, "true"},
	{`4.99 | default: 2.99`, 4.99},
	{`fruits | default: 2.99 | join`, "apples oranges peaches plums"},

	// array filters
	{`pages | map: 'category' | join`, "business celebrities lifestyle sports technology"},
	{`pages | map: 'category' | compact | join`, "business celebrities lifestyle sports technology"},
	{`"John, Paul, George, Ringo" | split: ", " | join: " and "`, "John and Paul and George and Ringo"},
	{`",John, Paul, George, Ringo" | split: ", " | join: " and "`, ",John and Paul and George and Ringo"},
	{`"John, Paul, George, Ringo," | split: ", " | join: " and "`, "John and Paul and George and Ringo,"},
	{`animals | sort | join: ", "`, "Sally Snake, giraffe, octopus, zebra"},
	{`sort_prop | sort: "weight" | inspect`, `[{"weight":null},{"weight":1},{"weight":3},{"weight":5}]`},
	{`fruits | reverse | join: ", "`, "plums, peaches, oranges, apples"},
	{`fruits | first`, "apples"},
	{`fruits | last`, "plums"},
	{`empty_array | first`, nil},
	{`empty_array | last`, nil},
	{`empty_array | last`, nil},
	{`dup_ints | uniq | join`, "1 2 3"},
	{`dup_strings | uniq | join`, "one two three"},
	{`dup_maps | uniq | map: "name" | join`, "m1 m2 m3"},
	{`mixed_case_array | sort_natural | join`, "a B c"},
	{`mixed_case_hash_values | sort_natural: 'key' | map: 'key' | join`, "a B c"},

	{`map_slice_has_nil | compact | join`, `a b`},
	{`map_slice_2 | first`, `b`},
	{`map_slice_2 | last`, `a`},
	{`map_slice_2 | join`, `b a`},
	{`map_slice_objs | map: "key" | join`, `a b`},
	{`map_slice_2 | reverse | join`, `a b`},
	{`map_slice_2 | sort | join`, `a b`},
	{`map_slice_dup | join`, `a a b`},
	{`map_slice_dup | uniq | join`, `a b`},

	// date filters
	{`article.published_at | date`, "Fri, Jul 17, 15"},
	{`article.published_at | date: "%a, %b %d, %y"`, "Fri, Jul 17, 15"},
	{`article.published_at | date: "%Y"`, "2015"},
	{`"2017-02-08 19:00:00 -05:00" | date`, "Wed, Feb 08, 17"},
	{`"2017-05-04 08:00:00 -04:00" | date: "%b %d, %Y"`, "May 04, 2017"},
	{`"2017-02-08 09:00:00" | date: "%H:%M"`, "09:00"},
	{`"2017-02-08 09:00:00" | date: "%-H:%M"`, "9:00"},
	{`"2017-02-08 09:00:00" | date: "%d/%m"`, "08/02"},
	{`"2017-02-08 09:00:00" | date: "%e/%m"`, " 8/02"},
	{`"2017-02-08 09:00:00" | date: "%-d/%-m"`, "8/2"},
	{`"March 14, 2016" | date: "%b %d, %y"`, "Mar 14, 16"},
	{`"2017-07-09" | date: "%d/%m"`, "09/07"},
	{`"2017-07-09" | date: "%e/%m"`, " 9/07"},
	{`"2017-07-09" | date: "%-d/%-m"`, "9/7"},
	{`ortto.example_date | date`, "Fri, Jul 17, 15"},
	{`ortto.not_existing_date | date`, ""},

	// sequence (array or string) filters
	{`"Ground control to Major Tom." | size`, 28},
	{`"apples, oranges, peaches, plums" | split: ", " | size`, 4},

	// string filters
	{`"Take my protein pills and put my helmet on" | replace: "my", "your"`, "Take your protein pills and put your helmet on"},
	{`"Take my protein pills and put my helmet on" | replace_first: "my", "your"`, "Take your protein pills and put my helmet on"},
	{`"/my/fancy/url" | append: ".html"`, "/my/fancy/url.html"},
	{`"website.com" | append: "/index.html"`, "website.com/index.html"},
	{`"title" | capitalize`, "Title"},
	{`"Élio Silva" | capitalize`, "Élio Silva"},
	{`"my great title" | capitalize`, "My great title"},
	{`"" | capitalize`, ""},
	{`"Parker Moore" | downcase`, "parker moore"},
	{`"Have you read 'James & the Giant Peach'?" | escape`, "Have you read &#39;James &amp; the Giant Peach&#39;?"},
	{`"1 < 2 & 3" | escape_once`, "1 &lt; 2 &amp; 3"},
	{`string_with_newlines | newline_to_br`, "<br />Hello<br />there<br />"},
	{`"1 &lt; 2 &amp; 3" | escape_once`, "1 &lt; 2 &amp; 3"},
	{`"apples, oranges, and bananas" | prepend: "Some fruit: "`, "Some fruit: apples, oranges, and bananas"},
	{`"I strained to see the train through the rain" | remove: "rain"`, "I sted to see the t through the "},
	{`"I strained to see the train through the rain" | remove_first: "rain"`, "I sted to see the train through the rain"},

	{`"Liquid" | slice: 0`, "L"},
	{`"Liquid" | slice: 2`, "q"},
	{`"Liquid" | slice: 2, 5`, "quid"},
	{`"Liquid" | slice: -3, 2`, "ui"},
	{`"Привет" | slice: -3, 2`, "ве"},

	{`"a/b/c" | split: '/' | join: '-'`, "a-b-c"},
	{`"a/b/" | split: '/' | join: '-'`, "a-b"},
	{`"a//c" | split: '/' | join: '-'`, "a--c"},
	{`"a//" | split: '/' | join: '-'`, "a"},
	{`"/b/c" | split: '/' | join: '-'`, "-b-c"},
	{`"/b/" | split: '/' | join: '-'`, "-b"},
	{`"//c" | split: '/' | join: '-'`, "--c"},
	{`"//" | split: '/' | join: '-'`, ""},
	{`"/" | split: '/' | join: '-'`, ""},
	{`"a.b" | split: '.' | join: '-'`, "a-b"},
	{`"a..b" | split: '.' | join: '-'`, "a--b"},
	{"'a.\t.b' | split: '.' | join: '-'", "a-\t-b"},
	{`"a b" | split: ' ' | join: '-'`, "a-b"},
	{`"a  b" | split: ' ' | join: '-'`, "a-b"},
	{"'a \t b' | split: ' ' | join: '-'", "a-b"},

	{`"Have <em>you</em> read <strong>Ulysses</strong>?" | strip_html`, "Have you read Ulysses?"},
	{`string_with_newlines | strip_newlines`, "Hellothere"},

	{`"Ground control to Major Tom." | truncate: 20`, "Ground control to..."},
	{`"Ground control to Major Tom." | truncate: 25, ", and so on"`, "Ground control, and so on"},
	{`"Ground control to Major Tom." | truncate: 20, ""`, "Ground control to Ma"},
	{`"Ground" | truncate: 20`, "Ground"},
	{`"Ground control to Major Tom." | truncatewords: 3`, "Ground control to..."},
	{`"Ground control to Major Tom." | truncatewords: 3, "--"`, "Ground control to--"},
	{`"Ground control to Major Tom." | truncatewords: 3, ""`, "Ground control to"},
	{`"Ground control" | truncatewords: 3, ""`, "Ground control"},
	{`"Ground" | truncatewords: 3, ""`, "Ground"},
	{`"  Ground" | truncatewords: 3, ""`, "  Ground"},
	{`"" | truncatewords: 3, ""`, ""},
	{`"  " | truncatewords: 3, ""`, "  "},

	{`"Parker Moore" | upcase`, "PARKER MOORE"},
	{`"          So much room for activities!          " | strip`, "So much room for activities!"},
	{`"          So much room for activities!          " | lstrip`, "So much room for activities!          "},
	{`"          So much room for activities!          " | rstrip`, "          So much room for activities!"},

	{`"%27Stop%21%27+said+Fred" | url_decode`, "'Stop!' said Fred"},
	{`"john@liquid.com" | url_encode`, "john%40liquid.com"},
	{`"Tetsuro Takara" | url_encode`, "Tetsuro+Takara"},

	// number filters
	{`"45" | to_number`, 45},
	{`-17 | abs`, 17},
	{`4 | abs`, 4},
	{`"-19.86" | abs`, 19.86},

	{`1.2 | ceil`, 2},
	{`2.0 | ceil`, 2},
	{`183.357 | ceil`, 184},
	{`"3.5" | ceil`, 4},

	{`1.2 | floor`, 1},
	{`2.0 | floor`, 2},
	{`183.357 | floor`, 183},

	{`4 | plus: 2`, 6},
	{`183.357 | plus: 12`, 195.357},

	{`4 | minus: 2`, 2},
	{`16 | minus: 4`, 12},
	{`183.357 | minus: 12`, 171.357},

	{`3 | times: 2`, 6},
	{`24 | times: 7`, 168},
	{`183.357 | times: 12`, 2200.284},

	{`3 | modulo: 2`, 1},
	{`24 | modulo: 7`, 3},
	// {`183.357 | modulo: 12 | `, 3.357}, // TODO test suit use inexact

	{`16 | divided_by: 4`, 4},
	{`5 | divided_by: 3`, 1},
	{`20 | divided_by: 7`, 2},
	{`20 | divided_by: 7.0`, 2.857142857142857},
	{`20 | divided_by: 's'`, nil},
	{`20 | divided_by: 0`, nil},

	{`1.2 | round`, 1},
	{`2.7 | round`, 3},
	{`183.357 | round: 2`, 183.36},

	// Jekyll extensions; added here for convenient testing
	// TODO add this just to the test environment
	{`map | inspect`, `{"a":1}`},
	{`1 | type`, `int`},
	{`"1" | type`, `string`},

	// Hash filters

	{`"Take my protein pills and put my helmet on" | md5`, "505a1a407670a93d9ef2cf34960002f9"},
	{`100 | md5`, "f899139df5e1059396431415e770c6dd"},
	{`100.01 | md5`, "e74f9831767648ecdd211c3f8cd85b86"},

	{`"Take my protein pills and put my helmet on" | sha1`, "07f3b4973325af9109399ead74f2180bcaefa4c0"},
	{`"" | sha1`, ""},
	{`100 | sha1`, "310b86e0b62b828562fc91c7be5380a992b2786a"},
	{`100.01 | sha1`, "2cf9b40e62dd0bff2c57d179bfc99674d25f3c33"},

	{`"Take my protein pills and put my helmet on" | sha256`, "b19c3d04c1b80ae9acd15227c0dde0cb6f5755995afa3c846a3473ac42de6f63"},
	{`"" | sha256`, ""},
	{`100 | sha256`, "ad57366865126e55649ecb23ae1d48887544976efea46a48eb5d85a6eeb4d306"},
	{`100.01 | sha256`, "4b46711a09b65af6dcbbc4caab38ab58e06d08eb75fbeb8e367fdd1ccc289fba"},

	{`"Take my protein pills and put my helmet on" | hmac: "key"`, "5b74077685d98d1e1d03cd289e2c2bfc"},
	{`"Take my protein pills and put my helmet on" | hmac: ""`, ""},
	{`"" | hmac: "key"`, ""},
	{`"" | hmac: 100`, ""},
	{`"" | hmac: 100.01`, ""},
	{`"Take my protein pills and put my helmet on" | hmac: 100`, "3494f6a7895d9e8084343e1020984ba6"},
	{`"Take my protein pills and put my helmet on" | hmac: 100.01`, "c1ef31ab6b3630ffb2e6842a600bf572"},
	{`"Only numeric and string keys are supported" | hmac: true`, ""},
	{`100 | hmac: "key"`, "f69388563202c10d4e0dc44646a3b937"},
	{`100 | hmac: 100`, "e459c4d00f32981388e5d0e797c8ac68"},
	{`100 | hmac: 100.01`, "f88e6d1df733b884b9748bbab83b3e68"},
	{`100.01 | hmac: "key"`, "41e66d9c6ca6e0b7b0470d9c03fef001"},
	{`100.01 | hmac: 100`, "7ac1da15168b6bf50c2975fa3198e84e"},
	{`100.01 | hmac: 100.01`, "bcd8551b5dbc26ed858752b9046dc654"},

	{`"Take my protein pills and put my helmet on" | hmac_sha1: "key"`, "fca4135e0bc4d4bcdccfd0bd98edc30d3d7ac629"},
	{`"Take my protein pills and put my helmet on" | hmac_sha1: ""`, ""},
	{`"" | hmac_sha1: "key"`, ""},
	{`"" | hmac_sha1: 100`, ""},
	{`"" | hmac_sha1: 100.01`, ""},
	{`"Take my protein pills and put my helmet on" | hmac_sha1: 100`, "595095014fab1b061a47cc1b7856b78bd78ad998"},
	{`"Take my protein pills and put my helmet on" | hmac_sha1: 100.01`, "3922875669b50f66373f1a21d91fd113f456b66c"},
	{`"Only numeric and string keys are supported" | hmac_sha1: true`, ""},
	{`100 | hmac_sha1: "key"`, "30385a0b6d754aee6a69093edd9d16accd57e26d"},
	{`100 | hmac_sha1: 100`, "56ba1ffa433eef7d9ebe9ef9fc464bdf2d68d7ed"},
	{`100 | hmac_sha1: 100.01`, "f962759dc0683e9aed4728d10cad6ade3c0f03ac"},
	{`100.01 | hmac_sha1: "key"`, "a3812ff53e8080fd42193b75d2245fe7ecb08df5"},
	{`100.01 | hmac_sha1: 100`, "877bfb3895f60525f123edec278d7dd915c6b2a6"},
	{`100.01 | hmac_sha1: 100.01`, "0efc1381dd2a001a0ba3db56f6e9456f3f4d73a8"},

	{`"Take my protein pills and put my helmet on" | hmac_sha256: "key"`, "111fce4b586c1c54804196bbc014e45005958fcaf5462fa206ad5856811686f5"},
	{`"Take my protein pills and put my helmet on" | hmac_sha256: ""`, ""},
	{`"" | hmac_sha256: "key"`, ""},
	{`"" | hmac_sha256: 100`, ""},
	{`"" | hmac_sha256: 100.01`, ""},
	{`"Take my protein pills and put my helmet on" | hmac_sha256: 100`, "c23af083390e2408faed6cf7d23f914425e9cab268050d5dc674f023bc8a8d6a"},
	{`"Take my protein pills and put my helmet on" | hmac_sha256: 100.01`, "9a19b23c1e55a2f570aad746844cb36f928d20ff4c837dca8fef0c2ef453cf63"},
	{`"Only numeric and string keys are supported" | hmac_sha256: true`, ""},
	{`100 | hmac_sha256: "key"`, "71d0fcbb40b55250039eb1f8bf363e280431f868af075355e6c9e44574f915d8"},
	{`100 | hmac_sha256: 100`, "f74a692209268d93c5a6ec227edfe17f7a70b28e049648f80238695798ffd407"},
	{`100 | hmac_sha256: 100.01`, "571751c3df688bc29af6e730c0c0d02ed4f1261fdfc9de2bf51a274106a5c6d4"},
	{`100.01 | hmac_sha256: "key"`, "b6c9391539ba7d250c9cbea6fb8aaaf278a5f858ad9206ae7ba6063ae17f2eb6"},
	{`100.01 | hmac_sha256: 100`, "7a48e1789185ab575a94579302ff9c4b57e58c70e40609f7a2a76469c9381d01"},
	{`100.01 | hmac_sha256: 100.01`, "bad95722cd8088216306962a575751a3a7251234f61504b33be224f9a9c2971c"},

	// at_least
	{`"10" | at_least: "20"`, 20},
	{`"10.5" | at_least: "20"`, 20},
	{`"10.5" | at_least: "20.5"`, 20.5},
	{`10 | at_least: 20`, 20},
	{`10.5 | at_least: 20`, 20},
	{`10.5 | at_least: 20.5`, 20.5},
	{`10 | at_least: "20"`, 20},
	{`10.5 | at_least: "20"`, 20},
	{`10.5 | at_least: "20.5"`, 20.5},
	{`"10" | at_least: 20`, 20},
	{`"10.5" | at_least: 20`, 20},
	{`"10.5" | at_least: 20.5`, 20.5},

	{`"20" | at_least: "10"`, 20},
	{`"20.5" | at_least: "10"`, 20.5},
	{`"20.5" | at_least: "10.5"`, 20.5},
	{`20 | at_least: 10`, 20},
	{`20.5 | at_least: 10`, 20.5},
	{`20.5 | at_least: 10.5`, 20.5},
	{`20 | at_least: "10"`, 20},
	{`20.5 | at_least: "10"`, 20.5},
	{`20.5 | at_least: "10.5"`, 20.5},
	{`"20" | at_least: 10`, 20},
	{`"20.5" | at_least: 10`, 20.5},
	{`"20.5" | at_least: 10.5`, 20.5},

	{`"0" | at_least: "0"`, 0},
	{`0 | at_least: "0"`, 0},
	{`"0" | at_least: 0`, 0},
	{`"0.0" | at_least: "0.0"`, 0},
	{`0.0 | at_least: "0.0"`, 0},
	{`"0.0" | at_least: 0.0`, 0},

	{`"" | at_least: 20`, ""},
	{`"" | at_least: "20"`, ""},
	{`"" | at_least: 20.5`, ""},
	{`"" | at_least: "20.5"`, ""},
	{`10 | at_least: ""`, ""},
	{`"10" | at_least: ""`, ""},
	{`"10.2" | at_least: ""`, ""},
	{`"10.2" | at_least: ""`, ""},

	// at_most
	{`"10" | at_most: "20"`, 10},
	{`"10.5" | at_most: "20"`, 10.5},
	{`"10.5" | at_most: "20.5"`, 10.5},
	{`10 | at_most: 20`, 10},
	{`10.5 | at_most: 20`, 10.5},
	{`10.5 | at_most: 20.5`, 10.5},
	{`10 | at_most: "20"`, 10},
	{`10.5 | at_most: "20"`, 10.5},
	{`10.5 | at_most: "20.5"`, 10.5},
	{`"10" | at_most: 20`, 10},
	{`"10.5" | at_most: 20`, 10.5},
	{`"10.5" | at_most: 20.5`, 10.5},

	{`"20" | at_most: "10"`, 10},
	{`"20.5" | at_most: "10"`, 10},
	{`"20.5" | at_most: "10.5"`, 10.5},
	{`20 | at_most: 10`, 10},
	{`20.5 | at_most: 10`, 10},
	{`20.5 | at_most: 10.5`, 10.5},
	{`20 | at_most: "10"`, 10},
	{`20.5 | at_most: "10"`, 10},
	{`20.5 | at_most: "10.5"`, 10.5},
	{`"20" | at_most: 10`, 10},
	{`"20.5" | at_most: 10`, 10},
	{`"20.5" | at_most: 10.5`, 10.5},

	{`"0" | at_most: "0"`, 0},
	{`0 | at_most: "0"`, 0},
	{`"0" | at_most: 0`, 0},
	{`"0.0" | at_most: "0.0"`, 0},
	{`0.0 | at_most: "0.0"`, 0},
	{`"0.0" | at_most: 0.0`, 0},

	{`"" | at_most: 20`, ""},
	{`"" | at_most: "20"`, ""},
	{`"" | at_most: 20.5`, ""},
	{`"" | at_most: "20.5"`, ""},
	{`10 | at_most: ""`, ""},
	{`"10" | at_most: ""`, ""},
	{`"10.2" | at_most: ""`, ""},
	{`"10.2" | at_most: ""`, ""},
}

var filterTestBindings = map[string]interface{}{
	"empty_array":     []interface{}{},
	"empty_map":       map[string]interface{}{},
	"empty_map_slice": yaml.MapSlice{},
	"map": map[string]interface{}{
		"a": 1,
	},
	"map_slice_2":       yaml.MapSlice{{Key: 1, Value: "b"}, {Key: 2, Value: "a"}},
	"map_slice_dup":     yaml.MapSlice{{Key: 1, Value: "a"}, {Key: 2, Value: "a"}, {Key: 3, Value: "b"}},
	"map_slice_has_nil": yaml.MapSlice{{Key: 1, Value: "a"}, {Key: 2, Value: nil}, {Key: 3, Value: "b"}},
	"map_slice_objs": yaml.MapSlice{
		{Key: 1, Value: map[string]interface{}{"key": "a"}},
		{Key: 2, Value: map[string]interface{}{"key": "b"}},
	},
	"mixed_case_array": []string{"c", "a", "B"},
	"mixed_case_hash_values": []map[string]interface{}{
		{"key": "c"},
		{"key": "a"},
		{"key": "B"},
	},
	"sort_prop": []map[string]interface{}{
		{"weight": 1},
		{"weight": 5},
		{"weight": 3},
		{"weight": nil},
	},
	"string_with_newlines": "\nHello\nthere\n",
	"dup_ints":             []int{1, 2, 1, 3},
	"dup_strings":          []string{"one", "two", "one", "three"},

	// for examples from liquid docs
	"animals": []string{"zebra", "octopus", "giraffe", "Sally Snake"},
	"fruits":  []string{"apples", "oranges", "peaches", "plums"},
	"article": map[string]interface{}{
		"published_at": timeMustParse("2015-07-17T15:04:05Z"),
	},
	"ortto": map[string]interface{}{
		"example_date": date.MustNewFromTime(timeMustParse("2015-07-17T15:04:05Z")),
	},
	"page": map[string]interface{}{
		"title": "Introduction",
	},
	"pages": []map[string]interface{}{
		{"name": "page 1", "category": "business"},
		{"name": "page 2", "category": "celebrities"},
		{"name": "page 3"},
		{"name": "page 4", "category": "lifestyle"},
		{"name": "page 5", "category": "sports"},
		{"name": "page 6"},
		{"name": "page 7", "category": "technology"},
	},
}

func TestFilters(t *testing.T) {
	require.NoError(t, os.Setenv("TZ", "America/New_York"))

	var (
		m1 = map[string]interface{}{"name": "m1"}
		m2 = map[string]interface{}{"name": "m2"}
		m3 = map[string]interface{}{"name": "m3"}
	)
	filterTestBindings["dup_maps"] = []interface{}{m1, m2, m1, m3}

	cfg := expressions.NewConfig()
	AddStandardFilters(&cfg)
	context := expressions.NewContext(filterTestBindings, cfg)

	for i, test := range filterTests {
		t.Run(fmt.Sprintf("%02d", i+1), func(t *testing.T) {
			actual, err := expressions.EvaluateString(test.in, context)
			require.NoErrorf(t, err, test.in)
			expected := test.expected
			switch v := actual.(type) {
			case int:
				actual = float64(v)
			}
			switch ex := expected.(type) {
			case int:
				expected = float64(ex)
			}
			require.Equalf(t, expected, actual, test.in)
		})
	}
}

func timeMustParse(s string) time.Time {
	t, err := time.Parse(time.RFC3339, s)
	if err != nil {
		panic(err)
	}
	return t
}
