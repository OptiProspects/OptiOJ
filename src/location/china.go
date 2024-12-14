package location

const (
	Province1  = "安徽"
	Province2  = "澳门"
	Province3  = "北京"
	Province4  = "重庆"
	Province5  = "福建"
	Province6  = "甘肃"
	Province7  = "广东"
	Province8  = "广西"
	Province9  = "贵州"
	Province10 = "海南"
	Province11 = "河北"
	Province12 = "河南"
	Province13 = "黑龙江"
	Province14 = "湖北"
	Province15 = "湖南"
	Province16 = "吉林"
	Province17 = "江苏"
	Province18 = "江西"
	Province19 = "辽宁"
	Province20 = "内蒙古自治区"
	Province21 = "宁夏"
	Province22 = "青海"
	Province23 = "山东"
	Province24 = "陕西"
	Province25 = "山西"
	Province26 = "上海"
	Province27 = "四川"
	Province28 = "台湾"
	Province29 = "天津"
	Province30 = "西藏自治区"
	Province31 = "香港"
	Province32 = "新疆"
	Province33 = "云南"
	Province34 = "浙江"
)

var Provinces = []string{
	Province1,
	Province2,
	Province3,
	Province4,
	Province5,
	Province6,
	Province7,
	Province8,
	Province9,
	Province10,
	Province11,
	Province12,
	Province13,
	Province14,
	Province15,
	Province16,
	Province17,
	Province18,
	Province19,
	Province20,
	Province21,
	Province22,
	Province23,
	Province24,
	Province25,
	Province26,
	Province27,
	Province28,
	Province29,
	Province30,
	Province31,
	Province32,
	Province33,
	Province34,
}

var CitiesMap = map[string][]string{
	"新疆":     {"阿克苏", "阿拉尔", "阿勒泰", "巴音郭楞", "博尔塔拉", "昌吉", "哈密", "和田", "喀什", "克拉玛依", "克孜勒苏柯尔克孜", "石河子", "塔城", "铁门关", "吐鲁番", "图木舒克", "五家渠", "乌鲁木齐", "伊犁哈萨克"},
	"贵州":     {"安顺", "毕节", "贵阳", "六盘水", "黔东南", "黔南", "黔西南", "铜仁", "遵义"},
	"湖北":     {"鄂州", "恩施", "黄冈", "黄石", "荆门", "荆州", "潜江", "神农架林区", "十堰", "随州", "天门", "武汉", "咸宁", "仙桃", "襄阳", "孝感", "宜昌"},
	"内蒙古自治区": {"阿拉善盟", "巴彦淖尔", "包头", "赤峰", "鄂尔多斯", "呼和浩特", "呼伦贝尔", "通辽", "乌海", "乌兰察布", "锡林郭勒盟", "兴安盟"},
	"西藏自治区":  {"阿里", "昌都", "拉萨", "林芝", "那曲", "日喀则", "山南"},
	"湖南":     {"常德", "长沙", "郴州", "衡阳", "怀化", "娄底", "邵阳", "湘潭", "湘西", "益阳", "永州", "岳阳", "张家界", "株洲"},
	"吉林":     {"白城", "白山", "长春", "吉林", "辽源", "四平", "松原", "通化", "延边"},
	"陕西":     {"安康", "宝鸡", "汉中", "商洛", "铜川", "渭南", "西安", "咸阳", "延安", "榆林"},
	"上海":     {},
	"广东":     {"潮州", "东莞", "佛山", "广州", "河源", "惠州", "江门", "揭阳", "茂名", "梅州", "清远", "汕头", "汕尾", "韶关", "深圳", "阳江", "云浮", "湛江", "肇庆", "中山", "珠海"},
	"海南":     {"白沙黎族自治县", "保亭黎族苗族自治县", "昌江黎族自治县", "澄迈县", "儋州", "定安县", "东方", "海口", "乐东黎族自治县", "临高县", "陵水黎族自治县", "琼海", "琼中黎族苗族自治县", "三沙", "三亚", "屯昌县", "万宁", "文昌", "五指山"},
	"河南":     {"安阳", "鹤壁", "济源", "焦作", "开封", "洛阳", "南阳", "平顶山", "濮阳", "三门峡", "商丘", "漯河", "新乡", "信阳", "许昌", "郑州", "周口", "驻马店"},
	"青海":     {"果洛", "海北", "海东", "海南", "海西", "黄南", "西宁", "玉树"},
	"台湾":     {},
	"香港":     {},
	"安徽":     {"安庆", "蚌埠", "亳州", "池州", "滁州", "阜阳", "合肥", "淮北", "淮南", "黄山", "六安", "马鞍山", "宿州", "铜陵", "芜湖", "宣城"},
	"澳门":     {},
	"江苏":     {"常州", "淮安", "连云港", "南京", "南通", "宿迁", "苏州", "泰州", "无锡", "徐州", "盐城", "扬州", "镇江"},
	"四川":     {"阿坝", "巴中", "成都", "达州", "德阳", "甘孜", "广安", "广元", "乐山", "凉山", "泸州", "眉山", "绵阳", "南充", "内江", "攀枝花", "遂宁", "雅安", "宜宾", "自贡", "资阳"},
	"重庆":     {},
	"广西":     {"百色", "北海", "崇左", "防城港", "贵港", "桂林", "河池", "贺州", "来宾", "柳州", "南宁", "钦州", "梧州", "玉林"},
	"山东":     {"滨州", "德州", "东营", "菏泽", "济南", "济宁", "莱芜", "聊城", "临沂", "青岛", "日照", "泰安", "潍坊", "威海", "烟台", "枣庄", "淄博"},
	"云南":     {"保山", "楚雄", "大理", "德宏", "迪庆", "红河", "昆明", "丽江", "临沧", "怒江", "普洱", "曲靖", "文山", "西双版纳", "玉溪", "昭通"},
	"江西":     {"抚州", "赣州", "吉安", "景德镇", "九江", "南昌", "萍乡", "上饶", "新余", "宜春", "鹰潭"},
	"山西":     {"大同", "晋城", "晋中", "临汾", "吕梁", "朔州", "太原", "忻州", "阳泉", "运城", "长治"},
	"天津":     {},
	"辽宁":     {"鞍山", "本溪", "大连", "丹东", "抚顺", "阜新", "葫芦岛", "锦州", "辽阳", "盘锦", "沈阳", "铁岭", "营口", "朝阳"},
	"浙江":     {"杭州", "湖州", "嘉兴", "金华", "丽水", "宁波", "衢州", "绍兴", "台州", "温州", "舟山"},
	"福建":     {"福州", "龙岩", "南平", "宁德", "莆田", "泉州", "三明", "厦门", "漳州"},
	"甘肃":     {"白银", "定西", "甘南", "嘉峪关", "金昌", "酒泉", "兰州", "临夏", "陇南", "平凉", "庆阳", "天水", "武威", "张掖"},
	"黑龙江":    {"大庆", "大兴安岭", "哈尔滨", "鹤岗", "黑河", "鸡西", "佳木斯", "牡丹江", "齐齐哈尔", "七台河", "双鸭山", "绥化", "伊春"},
	"北京":     {},
	"河北":     {"保定", "沧州", "承德", "定州", "邯郸", "衡水", "廊坊", "秦皇岛", "石家庄", "唐山", "辛集", "邢台", "张家口"},
	"宁夏":     {"固原", "石嘴山", "吴忠", "银川", "中卫"},
}

// GetCities 获取指定省份的城市列表
func GetCities(province string) []string {
	return CitiesMap[province]
}

// IsValidProvince 检查省份是否有效
func IsValidProvince(province string) bool {
	for _, p := range Provinces {
		if p == province {
			return true
		}
	}
	return false
}

// IsValidCity 检查城市是否属于指定省份
func IsValidCity(province, city string) bool {
	cities, ok := CitiesMap[province]
	if !ok {
		return false
	}
	for _, c := range cities {
		if c == city {
			return true
		}
	}
	return false
}