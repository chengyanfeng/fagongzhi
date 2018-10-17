var address = ""
var province = ""
var city = ""
var street = ""
var lng = ""
var lat = ""
var myselflastid=""
var returntotle = ""
var liebie = "-1";
var userid=""
var loginuserid=getCookie('userid')
var curNavIndex = 0;//首页0; 奶粉1; 面膜2; 图书3;
var lastId=""
var totalPage1=0
//初始化四个mescroll列表
var mescrollArr = new Array(4);//4个菜单所对应的4个mescroll对象
//模拟收索关键词
var curWord = '';
//初始化首页
mescrollArr[0] = initMescroll("mescroll0", "dataList0");
debugger
ifLogin()


/*alert(valueaa)*/
//初始化滚动组件
function initMescroll(mescrollId, clearEmptyId) {
    //创建MeScroll对象,内部已默认开启下拉刷新,自动执行up.callback,重置列表数据;
    var mescroll = new MeScroll(mescrollId, {
        down: {
            auto: false,//是否在初始化完毕之后自动执行下拉回调callback; 默认true; (注: down的callback默认调用 mescroll.resetUpScroll(); )
//					callback:function(mescroll) {
            //加载轮播数据
            //loadSwiper();
            //下拉刷新的回调,默认重置上拉加载列表为第一页(down的auto默认true,初始化Mescroll之后会自动执行到这里,而mescroll.resetUpScroll会触发up的callback)
//						mescroll.resetUpScroll();
//					}

            callback: downCallback

        },
        up: {
            auto:true, //初始化，上拉的时候默认加载第一页
            callback: getListData, //上拉回调,此处可简写; 相当于 callback: function (page) { getListData(page); }
            clearEmptyId: clearEmptyId, //1.下拉刷新时会自动先清空此列表,再加入数据; 2.无任何数据时会在此列表自动提示空
            isBounce: false, //此处禁止ios回弹,解析(务必认真阅读,特别是最后一点): http://www.mescroll.com/qa.html#q10
            noMoreSize: 3, //如果列表已无数据,可设置列表的总数量要大于半页才显示无更多数据;避免列表数据过少(比如只有一条数据),显示无更多数据会不好看
            empty: {
                //列表第一页无任何数据时,显示的空提示布局; 需配置warpId或clearEmptyId才生效;
                //warpId:null, //父布局的id; 如果此项有值,将不使用clearEmptyId的值;
                icon: "./static/images/mescroll-empty.png", //图标,默认null
                tip: "亲,没有您要找的信息~", //提示
//						btntext: "去逛逛 >", //按钮,默认""
//						btnClick: function(){//点击按钮的回调,默认null
//							alert("点击了按钮,具体逻辑自行实现");
//						}
            },
            toTop: { //配置回到顶部按钮
                //src : "../res/img/mescroll-totop.png", //默认滚动到1000px显示,可配置offset修改
                html: "<p>^<br/>顶部<p>", //标签内容,默认null; 如果同时设置了src,则优先取src
                offset: 500
            },
            lazyLoad: {
                use: false // 是否开启懒加载,默认false
            }
        }
    });


    return mescroll
}

/*下拉刷新的回调 */
function downCallback() {
    if (lng==""){
        alert("强烈建议您开启定位")
        //优先获取经纬度
        getjingweidu()
    }
    debugger
    //加载轮播数据..

    //分页id全部清空
    lastId=""
    myselflastid=""

    //加载列表数据
        debugger

        //重新设置当前的page的下拉加载
       mescrollArr[curNavIndex].resetUpScroll();
        /*mescrollArr[curNavIndex].setPageNum(0);*/






        /*//设置列表数据
        setListData(curPageData, curNavIndex);*/


}

/*切换列表*/
function changePage(i) {
    if (curNavIndex != i) {
        //更改列表条件
        $("#nav p").each(function (n, dom) {
            if (dom.getAttribute("i") == i) {
                dom.classList.add("active");
            } else {
                dom.classList.remove("active");
            }
        })
        //隐藏当前回到顶部按钮
        mescrollArr[curNavIndex].hideTopBtn();
        //取出菜单所对应的mescroll对象,如果未初始化则初始化
        if (mescrollArr[i] == null) {
            mescrollArr[i] = initMescroll("mescroll" + i, "dataList" + i);
        } else {
            //检查是否需要显示回到到顶按钮
            var curMescroll = mescrollArr[i];
            var curScrollTop = curMescroll.getScrollTop();
            if (curScrollTop >= curMescroll.optUp.toTop.offset) {
                curMescroll.showTopBtn();
            } else {
                curMescroll.hideTopBtn();
            }
        }
        //更新标记
        curNavIndex = i;
    }
}


/*联网加载列表数据  page = {num:1, size:10}; num:当前页 从1开始, size:每页数据条数 */
function getListData(page) {
        debugger
    var dataIndex = curNavIndex; //记录当前联网的nav下标,防止快速切换时,联网回来curNavIndex已经改变的情况;
    if(dataIndex==0){
        userid=""
    }else {
        if (loginuserid==""){
            userid=-1
        }else {
            userid=loginuserid
        }

    }
    //联网加载数据
    getListDataFromNet(curWord, liebie, dataIndex, page.num, page.size, function (curPageData, totalPage) {
        //联网成功的回调,隐藏下拉刷新和上拉加载的状态;
        //mescroll会根据传的参数,自动判断列表如果无任何数据,则提示空;列表无下一页数据,则提示无更多数据;
        console.log("dataIndex=" + dataIndex + "page.num=" + page.num + ", page.size=" + page.size + ", curPageData.length=" + curPageData.length);
        debugger
        //方法一(推荐): 后台接口有返回列表的总页数 totalPage
        mescrollArr[dataIndex].endByPage(curPageData.length, totalPage); //必传参数(当前页的数据个数, 总页数)
        //方法二(推荐): 后台接口有返回列表的总数据量 totalSize
        //mescroll.endBySize(curPageData.length, totalSize); //必传参数(当前页的数据个数, 总数据量)

        //方法三(推荐): 您有其他方式知道是否有下一页 hasNext
        //mescroll.endSuccess(curPageData.length, hasNext); //必传参数(当前页的数据个数, 是否有下一页true/false)

        //方法四 (不推荐),会存在一个小问题:比如列表共有20条数据,每页加载10条,共2页.如果只根据当前页的数据个数判断,则需翻到第三页才会知道无更多数据,如果传了hasNext,则翻到第二页即可显示无更多数据.
        //mescroll.endSuccess(curPageData.length);

        //提示:curPageData.length必传的原因:
        // 1.判断是否有下一页的首要依据: 当传的值小于page.size时,则一定会认为无更多数据.
        // 2.比传入的totalPage, totalSize, hasNext具有更高的判断优先级
        // 3.使配置的noMoreSize生效

        //设置列表数据,因为配置了emptyClearId,第一页会清空dataList的数据,所以setListData应该写在最后;
        debugger
        setListData(curPageData, dataIndex);
    }, function () {
        //联网失败的回调,隐藏下拉刷新和上拉加载的状态;
        mescrollArr[dataIndex].endErr();
    });
}

//热门搜索
$(".index-nav a").click(function () {
    debugger
    //把lastid置为空
    lastId=""
    curWord = $("#keyword").val(); //更新关键词
    liebie = $(this).attr("dd_name")
    //所有的变成未点击的样式
    $(".index-nav a").css("color","#4d525d")
    $(".index-nav a").css("text-decoration","none")

    //把当前字体变红
    $(this).css("text-decoration","underline")
    $(this).css("color","red")
    mescrollArr[0].resetUpScroll(); //重新搜索,重置列表数据
})

//搜索按钮
$("#search").click(function () {
    //把lastid置为空
    lastId=""
    debugger
    var word = $("#keyword").val();

        curWord = word; //更新关键词
        mescrollArr[0].resetUpScroll(); //重新搜索,重置列表数据

})
//小x号
function clearKey() {
    //清空
    $("#keyword").val("")
    curWord=""
    $("#keywordclear").css("display","none")
}

/*设置列表数据*/
function setListData(curPageData, dataIndex) {
    debugger
    console.log(curPageData)
    debugger
    if (loginuserid==0){
        $("#dataList2 .empty-tip").html("您还没有登陆哦")
    }else {
        if (curPageData.length==0){
            $("#dataList2 .empty-tip").html("您还没有发布信息哦")
        }
    }

    var listDom = document.getElementById("dataList" + dataIndex);

    for (var i = 0; i < curPageData.length; i++) {
        var pd = curPageData[i];
        var str = ""
        if (pd.liebie == "0") {
            str = '<p class="pd-name" style="display: inline-block">' + '拼车信息' + '</p>';
            if (pd.prize == "") {
                str += '<p class="pd-sold" style="display: inline-block;float:right">' + '费用:面议' + '</p> ';
            } else {
                str += '<p class="pd-sold" style="display: inline-block;float:right">' + '费用:' + pd.prize + '</p> ';
            }
            str += '<p class="">' + pd.local + '------->' + pd.destion + ' </p>';
        }
        if (pd.liebie == "1") {
            str = '<p class="pd-name" style="display: inline-block">' + '维修信息' + '</p>';
            if (pd.prize == "") {
                str += '<p class="pd-sold" style="display: inline-block;float:right">' + '红包:面议' + '</p> ';
            } else {
                str += '<p class="pd-sold" style="display: inline-block;float:right">' + '红包:' + pd.prize + '</p> ';
            }
        }
        if (pd.liebie == "2") {
            str = '<p class="pd-name" style="display: inline-block">' + '稍送信息' + '</p>';
            if (pd.prize == "") {
                str += '<p class="pd-sold" style="display: inline-block;float:right">' + '红包:面议' + '</p> ';
            } else {
                str += '<p class="pd-sold" style="display: inline-block;float:right">' + '红包:' + pd.prize + '</p> ';
            }
        }
        if (pd.liebie == "3") {
            str = '<p class="pd-name" style="display: inline-block">' + '借用信息' + '</p>';
            if (pd.prize == "") {
                str += '<p class="pd-sold" style="display: inline-block;float:right">' + '红包:面议' + '</p> ';
            } else {
                str += '<p class="pd-sold" style="display: inline-block;float:right">' + '红包:' + pd.prize + '</p> ';
            }

        }
        if (pd.liebie == "4") {
            str = '<p class="pd-name" style="display: inline-block">' + '二手物品信息' + '</p>';
            if (pd.prize == "") {
                str += '<p class="pd-sold" style="display: inline-block;float:right">' + '价格:面议' + '</p> ';
            } else {
                str += '<p class="pd-sold" style="display: inline-block;float:right">' + '价格:' + pd.prize + '</p> ';
            }
        }
        if (pd.liebie == "5") {
            str = '<p class="pd-name" style="display: inline-block">' + '帮助信息' + '</p>';
            if (pd.prize == "") {
                str += '<p class="pd-sold" style="display: inline-block;float:right">' + '红包:面议' + '</p> ';
            } else {
                str += '<p class="pd-sold" style="display: inline-block;float:right">' + '红包:' + pd.prize + '</p> ';
            }
        }


        str += '<p class="pd-price">' + pd.message + '</p>';
        if (pd.distance == "" || dataIndex == "2") {
            str += '<p class="pd-sold" style="display: inline-block">' + pd.city + pd.street + '</p> ';

        } else {
            str += '<p class="pd-sold" style="display: inline-block">距离您的位置' + pd.distance + 'km</p> ';

        }

        str += '<p class="pd-sold" style="display: inline-block;float:right">' + formatDateTime(pd.time) + '</p> ';

        var liDom = document.createElement("li");
        liDom.innerHTML = str;
        listDom.appendChild(liDom);

    }
}

/*联网加载列表数据
 在您的实际项目中,请参考官方写法: http://www.mescroll.com/api.html#tagUpCallback
 请忽略getListDataFromNet的逻辑,这里仅仅是在本地模拟分页数据,本地演示用
 实际项目以您服务器接口返回的数据为准,无需本地处理分页.
 * */
function getListDataFromNet(curWord, liebie, curNavIndex, pageNum, pageSize, successCallback, errorCallback) {
    //延时一秒,模拟联网
    setTimeout(function () {
        var param = "?curWord=" + curWord + "&lng=" + lng + "&lat=" + lat + "&liebie=" + liebie + "&userid=" + userid+"&myselflastid="+myselflastid+"&lastId="+lastId
        $.ajax({
            async:false,
            type: 'GET',
            url: 'https://chengyanfeng.natapp4.cc/TestGetPersion' + param,
//		                url: '../res/pdlist1.json?num='+pageNum+"&size="+pageSize+"&word="+curWord,
            dataType: 'json',
            success: function (dataAll) {
                if($("#keyword").val().length>0){
                    $("#keywordclear").css("display","inline")
                }

                debugger
                //判断是否获取经纬度
                if (lat==""){
                    if(lat==""){
                        getjingweidu()
                    }
                }

                var data = [];
                var listData = [];
                if (curNavIndex==0){
                    returntotle = dataAll.UserList.length
                    for (var i = 0; i < dataAll.UserList.length; i++) {

                        data.push(dataAll.UserList[i]);
                        listData.push(dataAll.UserList[i]);
                    }
                    if (returntotle>0){
                        //最后一个id
                        lastId= dataAll.UserList[returntotle-1].ID
                    }

                }else {
                    returntotle = dataAll.UserListLogin.length
                    for (var i = 0; i < dataAll.UserListLogin.length; i++) {

                        data.push(dataAll.UserListLogin[i]);
                        listData.push(dataAll.UserListLogin[i]);
                    }
                    if (returntotle>0){
                        //最后一个id
                        myselflastid= dataAll.UserListLogin[returntotle-1].ID
                    }
                }







                if (totalPage1==0){
                    totalPage1=dataAll.totalPage
                }
                debugger
                if (pageNum==1){
                    //查看是否第一个tab主页的滑动
                    if ($(".tab.active .data-list").attr("id") == "dataList0") {
                        //显示提示
                        $("#downloadTip").css("opacity", "1");
                        debugger
                        $("#downloadTip").html("更新" + returntotle + "条信息");
                        setTimeout(function () {
                            $("#downloadTip").css("opacity", "0");
                        }, 2000);
                    }


                    //联网成功的回调,隐藏下拉刷新的状态
                    mescrollArr[curNavIndex].endSuccess();

                        //获取当前的tab, //删除老数据
                        $(".tab.active .data-list").html("")




                }
                successCallback(listData, totalPage1);

            },
            error: errorCallback
        });
    }, 1000)
}


/*初始化菜单*/
$("#nav a").click(function () {
    var i = Number($(this).attr("i"));
    divshow(i)
    //如果是发布界面，那么不请求，任何数据，也不加载mescroll
    if (i == 1) {
        return
    } else {
        changePage(i)
    }

})


//div 显示与隐藏
function divshow(i) {
    if(i=="0"){
        userid=""
    }
    debugger
    $(".mescroll").each(function (n, dom) {
        if (dom.getAttribute("id") == "mescroll" + i) {
            //显示整个tab div 层
            $("#tab" + i).addClass("active")
            $("#tab" + i).removeClass("hidden")
            //显示滑动区域
            dom.classList.remove("hidden");
            dom.classList.add("active");
        } else {
            $("#tab" + n).removeClass("active")
            $("#tab" + n).addClass("hidden")
            dom.classList.remove("active");
            dom.classList.add("hidden");
        }
    })
}

/*切换列表*/
function changePage(i) {
    if (curNavIndex != i) {
        //更改列表条件
        $("#nav a").each(function (n, dom) {
            if (dom.getAttribute("i") == i) {
                dom.classList.add("active");
            } else {
                dom.classList.remove("active");
            }
        })
        //隐藏当前回到顶部按钮
        mescrollArr[curNavIndex].hideTopBtn();
        //取出菜单所对应的mescroll对象,如果未初始化则初始化
        if (mescrollArr[i] == null) {
            mescrollArr[i] = initMescroll("mescroll" + i, "dataList" + i);
        } else {
            //检查是否需要显示回到到顶按钮
            var curMescroll = mescrollArr[i];
            var curScrollTop = curMescroll.getScrollTop();
            if (curScrollTop >= curMescroll.optUp.toTop.offset) {
                curMescroll.showTopBtn();
            } else {
                curMescroll.hideTopBtn();
            }
        }
        //更新标记
        curNavIndex = i;
    }
}

//发布信息，select 筛选
function selectDiv() {
    i = $("#select  option:selected").val()
    debugger
    $("#mescroll1 .keyword").each(function (n, dom) {
        debugger
        if (i == n) {
            dom.classList.add("active")
            dom.classList.remove("hidden")
        } else {
            dom.classList.add("hidden")
            dom.classList.remove("active")
        }


    });
}
//发布

function UpMessage() {


    n = $("#select  option:selected").val()
    var formdata = new FormData()
    if (n == 0) {
        var local = $("#local").val()
        var destion = $("#destion").val()
        if (local == "" || destion == "") {
            alert("请填写乘车信息")
            return
        }
        formdata.append("local", local)
        formdata.append("destion", destion)
    } else {
        var keyword = $(".keyword.active :input").val()
        if (keyword == "") {
            alert("请填写关键词信息")
            return
        }
        formdata.append("keyword", keyword)
    }
    if ($("#message").val() == "") {
        alert("请填写信息")
        return
    }
    if ($("#message").val().length>200) {
        alert("信息超出字数限制")
        return
    }
    if (lat==""){
        alert("强烈建议您同意获取位置，周边的人可以根据位置获取您的信息")
        getjingweidu()
    }
    formdata.append("leibie", n)
    formdata.append("Exhour", $("#input_from").val())
    formdata.append("message", $("#message").val())
    formdata.append("prize", $("#thankpackage").val())
    formdata.append("province", province)
    formdata.append("city", city)
    formdata.append("street", street)
    formdata.append("address", address)
    formdata.append("lng", lng)
    formdata.append("lat", lat)
    if (loginuserid.length>0){
        formdata.append("userid",loginuserid)
    }else {
        alert("您还没有登陆，发布的信息将过期删除")
    }
    UpMessageAjax(formdata)

}

function UpMessageAjax(formdata) {
    $.ajax({
        type: 'post',
        url: 'https://chengyanfeng.natapp4.cc/upmessage',
        /* url: 'https://127.0.0.1/upmessage',*/
        dataType: 'json',
        data: formdata,
        processData: false,
        contentType: false,
        success: function (dataAll) {
            alert("信息上传成功")
            divshow("0")
            //清空列表
            $("#keyword").val("")
            curWord ="" ; //更新关键词
            liebie = "-1"
            //lastid=0
            lastId=""
            debugger
            //刷新列表
            curNavIndex=0
            mescrollArr[0].resetUpScroll();

            //清空发布信息
            $(".needclear").val("")

            //底部颜色
            upmessagebottom()
        },

    });
}


//********************************************时间插件**********************************************/
var $input = $('.datepicker').pickatime({
    format: 'H',
    clear: "",
    formatSubmit: 'H',
    interval: 60,
})
var picker = $input.pickatime('picker')
//****************************************************************************************************/

//********************************************高德地图接口**********************************************/
//获取高德地理对象
$(function () {

})
//******************************************************************************************/

//********************************************搜索关键字**********************************************/

//********************************************时间戳**********************************************/

function formatDateTime(inputTime) {
    var date = new Date(inputTime * 1000);
    var y = date.getFullYear();
    var m = date.getMonth() + 1;
    m = m < 10 ? ('0' + m) : m;
    var d = date.getDate();
    d = d < 10 ? ('0' + d) : d;
    var h = date.getHours();
    h = h < 10 ? ('0' + h) : h;
    var minute = date.getMinutes();
    var second = date.getSeconds();
    minute = minute < 10 ? ('0' + minute) : minute;
    second = second < 10 ? ('0' + second) : second;
    return y + '-' + m + '-' + d + ' ' + h + ':' + minute + ':' + second;
};

function showdiv() {
    document.getElementById("bg").style.display = "block";
    document.getElementById("zhuche").style.display = "block";
}

function showdenglu() {
    document.getElementById("bg").style.display = "block";
    document.getElementById("denglu").style.display = "block";

}
function showadvice(){
    document.getElementById("bg").style.display = "block";
    document.getElementById("showadvice").style.display = "block";

}

function sharefriend() {
    document.getElementById("bg").style.display = "block";
    document.getElementById("sharefriend").style.display = "block";
}
function downapp() {
    document.getElementById("bg").style.display = "block";
    document.getElementById("downapp").style.display = "block";
}

function hidediv() {
    document.getElementById("bg").style.display = 'none';
    document.getElementById("zhuche").style.display = 'none';
    document.getElementById("denglu").style.display = 'none';
    document.getElementById("getpassword").style.display = 'none';
    document.getElementById("showadvice").style.display = 'none';
    document.getElementById("sharefriend").style.display = "none";
    document.getElementById("downapp").style.display = "none";


    $("#getcode").css("display", "none")
    $("#getidentifying").css("display", "inline")
    $("#code").css("display", "none")
    $("#identifying").css("display", "inline")
    $(".needclear").val("")
}


function showgetpassword() {
    document.getElementById("bg").style.display = 'block';
    document.getElementById("denglu").style.display = 'none';
    document.getElementById("getpassword").style.display = 'block';
    $("#getcode").css("display", "none")
    $("#getidentifying").css("display", "inline")
}

//获取验证码
function identifying(id,identifyingid,codeid) {
    if (sendCode($("#"+id).val())==false) {
        return
    }
    $("#"+identifyingid).css("display", "none")
    $("#"+codeid).css("display", "inline")

}


//注册
function register() {
    number=$("#renumber").val()
    password=$("#repassword").val()
    code=$("#code").val()
    if (number==""){
        alert("电话号码为空")
        return
    }
    if (password==""){
        alert("请输入密码")
        return
    }
    if (code==""){
        alert("请输入验证码")
        return
    }
    var formdata = new FormData()
    formdata.append("phoneNumber",number)
    formdata.append("password",password)
    formdata.append("code",code)
    Ajax(formdata,"Register")

}


//登陆
function Login() {
    number=$("#number").val()
    password=$("#password").val()
    if (number==""){
        alert("电话号码为空")
        return
    }
    if (password==""){
        alert("请输入密码")
        return
    }
    var formdata = new FormData()
    formdata.append("phoneNumber",number)
    formdata.append("password",password)
    Ajax(formdata,"Login")

}

function getidentifying() {
    $("#getidentifying").css("display", "none")
    $("#getcode").css("display", "inline")
}
function sleep(n) { //n表示的毫秒数
    var start = new Date().getTime();
    while (true) if (new Date().getTime() - start > n) break;
}


//发送验证码的接口
function sendCode(number) {
    if (number==""){
        alert("请输入手机号")
        return false
    }else {
        debugger
        var formdata = new FormData()
        formdata.append("phoneNumber",number)
        Ajax(formdata,"SendCode")
    }
}



//发送建议的接口
function sendAdvice() {
        var formdata = new FormData()
         upadvice=$("#upadvice").val()
    if (upadvice.length<1){
            alert("请输入信息")
            return
    }
        formdata.append("upadvice",upadvice)
        Ajax(formdata,"UpAdvice")

}





//公用登陆ajax
function Ajax(formdata,method) {
    debugger
    $.ajax({
        type: 'post',
        url: 'https://chengyanfeng.natapp4.cc/'+method,
        /* url: 'https://127.0.0.1/upmessage',*/
        dataType: 'json',
        data: formdata,
        processData: false,
        contentType: false,
        success: function (dataAll) {
            if (method=="SendCode"){
                alert("验证码已经发送到您的手机，请在两分钟内操作")
            }
            if (method=="Login"){
                if (dataAll.Message=="ok"){
                    loginuserid=dataAll.Token
                    setCookie('userid',dataAll.Token)

                    //登陆成功，css 样式控制
                    hidediv()
                    alert("登陆成功")
                    //界面转换
                    hivLoginReg()
                    //刷新信息
                    mescrollArr[2].resetUpScroll();

                }else {
                    alert("密码错误，请重新输入密码")
                }


            }
            if (method=="Register"){
                if (dataAll.Count==0){
                    alert("已经注册成功请登录")
                    hidediv()
                }else{
                    alert("该手机号码已经注册过了")
                }

            }
            if (method=="ResTPassWord"){
                alert(dataAll.Message)
                    hidediv()

            }
            if (method=="UpAdvice"){
                alert(dataAll.Message)
                hidediv()

            }
        },

    });
}
//登陆成功后显示
function hivLoginReg() {

    //隐藏登陆与注册
    $("#register").css("display","none")
    $("#land").css("display","none")
    //显示欢迎和退出
    $("#welcome").css("display","")
    $("#quit").css("display","")

}
//未登陆显示
function showLoginReg() {

    //隐藏登陆与注册
    $("#register").css("display","")
    $("#land").css("display","")
    //显示欢迎和退出
    $("#welcome").css("display","none")
    $("#quit").css("display","none")

}


function setCookie(login,value)
{
    debugger
    var Days = 30;
    var exp = new Date();
    exp.setTime(exp.getTime() + Days*24*60*60*1000);
    document.cookie = login + "="+ escape (value) + ";expires=" + exp.toGMTString();
}

function getCookie(name)
{
    debugger
    var arr,reg=new RegExp("(^| )"+name+"=([^;]*)(;|$)");

    if(arr=document.cookie.match(reg))

        return unescape(arr[2]);
    else
        return "";
}

function setPassword() {
    number=$("#getnumber").val()
    password=$("#getrepassword").val()
    code=$("#getcode").val()
    if (number==""){
        alert("电话号码为空")
        return
    }
    if (password==""){
        alert("请输入密码")
        return
    }
    if (code==""){
        alert("请输入验证码")
        return
    }
    var formdata = new FormData()
    formdata.append("phoneNumber",number)
    formdata.append("password",password)
    formdata.append("code",code)
    Ajax(formdata,"ResTPassWord")
}

//退出
function quit() {
    debugger
    setCookie('userid',"")
    userid=""
    loginuserid=""
    myselflastid=""
    alert("您已经成功退出")
    mescrollArr[2].resetUpScroll();
    showLoginReg();

}






//判断是否登陆
function ifLogin() {
    if (loginuserid==""){
/*
        alert("显示登陆")
*/
        //显示登陆
        showLoginReg()
    }else {
        hivLoginReg()
    }
}

function upmessagebottom(){
    $(".nav a").css("color","#4d525d")
    $(".nav a").css("text-decoration","none")
    //第一个a 标签为红色
    $(".nav a:first").css("color","red")
    $(".nav a:first").css("text-decoration","underline")
}


//
$(".nav a").click(function () {
    $(".nav a").css("color","#4d525d")
    $(".nav a").css("text-decoration","none")
    //当前为红色
    $(this).css("color","red")
    $(this).css("text-decoration","underline")
});




/*经纬度*/
function getjingweidu(){
    var mapObj = new AMap.Map('iCenter');
    mapObj.plugin('AMap.Geolocation', function () {
        geolocation = new AMap.Geolocation({
            enableHighAccuracy: true, // 是否使用高精度定位，默认:true
            timeout: 10000,           // 超过10秒后停止定位，默认：无穷大
            maximumAge: 0,            // 定位结果缓存0毫秒，默认：0
            convert: true,            // 自动偏移坐标，偏移后的坐标为高德坐标，默认：true
            showButton: true,         // 显示定位按钮，默认：true
            buttonPosition: 'LB',     // 定位按钮停靠位置，默认：'LB'，左下角
            buttonOffset: new AMap.Pixel(10, 20), // 定位按钮与设置的停靠位置的偏移量，默认：Pixel(10, 20)
            showMarker: true,         // 定位成功后在定位到的位置显示点标记，默认：true
            showCircle: true,         // 定位成功后用圆圈表示定位精度范围，默认：true
            panToLocation: true,      // 定位成功后将定位到的位置作为地图中心点，默认：true
            zoomToAccuracy: true       // 定位成功后调整地图视野范围使定位位置及精度范围视野内可见，默认：false
        });
        mapObj.addControl(geolocation);
        geolocation.getCurrentPosition();
        AMap.event.addListener(geolocation, 'complete', onComplete); // 返回定位信息
        AMap.event.addListener(geolocation, 'error', onError);       // 返回定位出错信息
    });
}
//获取高德地理信息
function onComplete(obj) {
    var res = '经纬度：' + obj.position +
        '\n精度范围：' + obj.accuracy +
        '米\n定位结果的来源：' + obj.location_type +
        '\n状态信息：' + obj.info +
        '\n地址：' + obj.formattedAddress +
        '\n地址信息：' + JSON.stringify(obj.addressComponent, null, 4);
    debugger

          /*  alert(res);*/

    console.log(obj)
    console.log(obj)
    if (obj.addressComponent.city == "") {
        debugger
        //城市为空，省为城市

        if (obj.addressComponent.province == "") {
            alert("获取地理位置失败")
        } else {
            debugger
            province=obj.addressComponent.province
            city=obj.addressComponent.province
            street=obj.addressComponent.district+obj.addressComponent.street
            address=obj.formattedAddress
            lat=obj.position.lat
            lng=obj.position.lng

            if (obj.addressComponent.province.length > 3) {
                province = obj.addressComponent.city.substring(start, [end])
                $("#location").text(province)
            } else {
                $("#location").text(obj.addressComponent.province)
            }

        }
    } else {
        province=obj.addressComponent.province
        city=obj.addressComponent.city
        street=obj.addressComponent.district+obj.addressComponent.street
        address=obj.formattedAddress
        lat   =obj.position.lat

        lng =obj.position.lng


        if (obj.addressComponent.city.length > 3) {
            city = obj.addressComponent.city.substring(start, [end])
            $("#location").text(city)
        } else {
            $("#location").text(obj.addressComponent.city)
        }

    }

}

function onError(obj) {
    /*alert(obj.info + '--' + obj.message);*/
    alert("您已经拒绝授权定位，强烈建议您去设置里开启定位")
    console.log(obj);
}

//获取高德两点之间距离
var p1 = [116.434027, 39.941037];
var p2 = [116.461665, 39.941564];
// 返回 p1 到 p2 间的地面距离，单位：米
var dis = AMap.GeometryUtil.distance(p1, p2);
/*alert(dis)*/


//复制
function myCopy(){

    var ele = document.getElementById("copyUrl");//ele是要复制的元素的对象

    ele.focus();

    // ele.select();
    ele.setSelectionRange(0, ele.value.length);



    if(document.execCommand('copy', false, null)){
        //success info
        alert("复制成功")
    } else{
        //fail info
        alert("复制失败")
    }

}



