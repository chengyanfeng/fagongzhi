$(function(){

    var curNavIndex=0;//首页0; 奶粉1; 面膜2; 图书3;
    //初始化四个mescroll列表
    var mescrollArr=new Array(4);//4个菜单所对应的4个mescroll对象

    debugger
    //初始化首页
    mescrollArr[0]=initMescroll("mescroll0", "dataList0");


    //初始化滚动组件
    function initMescroll(mescrollId,clearEmptyId) {
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
            },
            up: {
                callback: getListData, //上拉回调,此处可简写; 相当于 callback: function (page) { getListData(page); }
                clearEmptyId: clearEmptyId, //1.下拉刷新时会自动先清空此列表,再加入数据; 2.无任何数据时会在此列表自动提示空
                isBounce: false, //此处禁止ios回弹,解析(务必认真阅读,特别是最后一点): http://www.mescroll.com/qa.html#q10
                noMoreSize: 3, //如果列表已无数据,可设置列表的总数量要大于半页才显示无更多数据;避免列表数据过少(比如只有一条数据),显示无更多数据会不好看
                empty: {
                    //列表第一页无任何数据时,显示的空提示布局; 需配置warpId或clearEmptyId才生效;
                    //warpId:null, //父布局的id; 如果此项有值,将不使用clearEmptyId的值;
                    icon: "../res/img/mescroll-empty.png", //图标,默认null
                    tip: "亲,没有您要找的商品~", //提示
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
                    use: true // 是否开启懒加载,默认false
                }
            }
        });
        return mescroll
    }


    /*切换列表*/
    function changePage(i) {
        if(curNavIndex!=i) {
            //更改列表条件
            $("#nav p").each(function(n,dom){
                if (dom.getAttribute("i")==i) {
                    dom.classList.add("active");
                } else{
                    dom.classList.remove("active");
                }
            })
            //隐藏当前回到顶部按钮
            mescrollArr[curNavIndex].hideTopBtn();
            //取出菜单所对应的mescroll对象,如果未初始化则初始化
            if(mescrollArr[i]==null){
                mescrollArr[i]=initMescroll("mescroll"+i, "dataList"+i);
            }else{
                //检查是否需要显示回到到顶按钮
                var curMescroll=mescrollArr[i];
                var curScrollTop=curMescroll.getScrollTop();
                if(curScrollTop>=curMescroll.optUp.toTop.offset){
                    curMescroll.showTopBtn();
                }else{
                    curMescroll.hideTopBtn();
                }
            }
            //更新标记
            curNavIndex=i;
        }
    }



    //模拟收索关键词
    var curWord='wokao';

    //热门搜索
    $(".hot-words li").click(function() {
        curWord=this.innerText; //更新关键词
        mescroll.resetUpScroll(); //重新搜索,重置列表数据
    })

    //搜索按钮
    $("#search").click(function(){
        var word=$("#keyword").val();
        if(word){
            curWord=word; //更新关键词
            mescroll.resetUpScroll(); //重新搜索,重置列表数据
        }
    })

    /*联网加载列表数据  page = {num:1, size:10}; num:当前页 从1开始, size:每页数据条数 */
    function getListData(page){

        var dataIndex=curNavIndex; //记录当前联网的nav下标,防止快速切换时,联网回来curNavIndex已经改变的情况;
        //联网加载数据
        getListDataFromNet(curWord,dataIndex, page.num, page.size, function(curPageData,totalPage){
            //联网成功的回调,隐藏下拉刷新和上拉加载的状态;
            //mescroll会根据传的参数,自动判断列表如果无任何数据,则提示空;列表无下一页数据,则提示无更多数据;
            console.log("dataIndex="+dataIndex+"page.num="+page.num+", page.size="+page.size+", curPageData.length="+curPageData.length);

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
            setListData(curPageData,dataIndex);
        }, function(){
            //联网失败的回调,隐藏下拉刷新和上拉加载的状态;
            mescrollArr[dataIndex].endErr();
        });
    }

    /*设置列表数据*/
    function setListData(curPageData,dataIndex){
        var listDom=document.getElementById("dataList"+dataIndex);
        for (var i = 0; i < curPageData.length; i++) {
            var pd=curPageData[i];
            var str='<p class="pd-name">'+pd.pdName+'</p>';
            str+='<p class="pd-price">'+pd.pdPrice+' 元</p>';
            str+='<p class="pd-sold">已售'+pd.pdSold+'件</p>';

            var liDom=document.createElement("li");
            liDom.innerHTML=str;
            listDom.appendChild(liDom);
        }
    }

    /*联网加载列表数据
     在您的实际项目中,请参考官方写法: http://www.mescroll.com/api.html#tagUpCallback
     请忽略getListDataFromNet的逻辑,这里仅仅是在本地模拟分页数据,本地演示用
     实际项目以您服务器接口返回的数据为准,无需本地处理分页.
     * */
    function getListDataFromNet(curWord,curNavIndex, pageNum,pageSize,successCallback,errorCallback) {
        //延时一秒,模拟联网
        setTimeout(function () {
            $.ajax({
                type: 'GET',
                url: 'http://chengyanfeng.natapp1.cc',
//		                url: '../res/pdlist1.json?num='+pageNum+"&size="+pageSize+"&word="+curWord,
                dataType: 'json',
                success: function(dataAll){
                    debugger
                    //模拟服务器接口的搜索
                    var data=[];
                    var listData=[];
                    for (var i = 0; i < dataAll.UserList.length; i++) {
                        if (dataAll.UserList[i].pdName.indexOf(curWord)!=-1) {
                            data.push(dataAll.UserList[i]);
                            listData.push(dataAll.UserList[i]);
                        }
                    }
                    //模拟服务器接口的分页

//							for (var i = (pageNum-1)*pageSize; i < pageNum*pageSize; i++) {
//			            		if(i==data.length) break;

//			            	}

                    successCallback(listData,dataAll.totalPage);
                },
                error: errorCallback
            });
        },500)
    }


    /*初始化菜单*/
    $("#nav a").click(function(){
        var i=Number($(this).attr("i"));
        divshow(i)
        //如果是发布界面，那么不请求，任何数据，也不加载mescroll
        if (i==1){
            return
            }else {
            changePage(i)
        }

    })
    //div 显示与隐藏
    function divshow(i){
        debugger
        $(".mescroll").each(function(n,dom){
            if (dom.getAttribute("id")=="mescroll"+i) {
                //显示整个tab div 层
                $("#tab"+i).addClass("active")
                $("#tab"+i).removeClass("hidden")
                //显示滑动区域
                dom.classList.remove("hidden");
                dom.classList.add("active");
            } else{
                $("#tab"+n).removeClass("active")
                $("#tab"+n).addClass("hidden")
                dom.classList.remove("active");
                dom.classList.add("hidden");
            }
        })
    }
    /*切换列表*/
    function changePage(i) {
        if(curNavIndex!=i) {
            //更改列表条件
            $("#nav a").each(function(n,dom){
                if (dom.getAttribute("i")==i) {
                    dom.classList.add("active");
                } else{
                    dom.classList.remove("active");
                }
            })
            //隐藏当前回到顶部按钮
            mescrollArr[curNavIndex].hideTopBtn();
            //取出菜单所对应的mescroll对象,如果未初始化则初始化
            if(mescrollArr[i]==null){
                mescrollArr[i]=initMescroll("mescroll"+i, "dataList"+i);
            }else{
                //检查是否需要显示回到到顶按钮
                var curMescroll=mescrollArr[i];
                var curScrollTop=curMescroll.getScrollTop();
                if(curScrollTop>=curMescroll.optUp.toTop.offset){
                    curMescroll.showTopBtn();
                }else{
                    curMescroll.hideTopBtn();
                }
            }
            //更新标记
            curNavIndex=i;
        }
    }
});

