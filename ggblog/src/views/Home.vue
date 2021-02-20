<template>
    <div class="home">
        <banner isHome="true" ></banner>
        <div class="site-content animate">
            <!--通知栏-->
            <div class="notify">
                <div class="search-result" v-if="hideSlogan">
                    <span v-if="searchWords">搜索结果："{{searchWords}}" 相关文章</span>
                    <span v-else-if="category">分类 "{{category}}" 相关文章</span>
                </div>
                <quote v-else>{{notice}}</quote>
            </div>

            <!--焦点图-->
            <div class="top-feature" v-if="!hideSlogan">
                <section-title>
                    <div style="display: flex;align-items: flex-end;">聚焦<small-ico></small-ico></div>
                </section-title>
                <div class="feature-content">
                    <div class="feature-item" v-for="item in features" :key="item.title">
                        <Feature :data="item"></Feature>
                    </div>
                </div>
            </div>
            <!--文章列表-->
            <main class="site-main" :class="{'search':hideSlogan}">
                <section-title v-if="!hideSlogan">推荐</section-title>
                <template v-for="item in postList">
                    <post :post="item" :key="item.id"></post>
                </template>
            </main>

            <!--加载更多-->
            <div class="more" v-show="hasNextPage">
                <div class="more-btn" @click="loadMore">More</div>
            </div>
        </div>
    </div>
</template>

<script>
    import Banner from '@/components/banner'
    import Feature from '@/components/feature'
    import sectionTitle from '@/components/section-title'
    import Post from '@/components/post'
    import SmallIco from '@/components/small-ico'
    import Quote from '@/components/quote'
    import {fetchFocus, fetchList,setting} from '../api'

    export default {
        name: 'Home',
        props: ['cate', 'words'],
        data() {
            return {
                features: [],
                postList: [],
                currPage: 0,
                hasNextPage: false
            }
        },
        components: {
            Banner,
            Feature,
            sectionTitle,
            Post,
            SmallIco,
            Quote
        },
        computed: {
            searchWords() {
                return this.$route.params.words
            },
            category() {
                return this.$route.params.cate
            },
            category_id() {
                return this.$route.params.cate_id
            },
            hideSlogan() {
                return this.category || this.searchWords
            },
            notice() {
                return this.$store.getters.notice
            }
        },
        watch: {
           '$route' (to,from) {
                      //对路由变化做出响应
                    
                    this.postList =   []
                    this.currPage = 0
                    this.hasNextPage = false
                        this.fetchFocus();
                    this.fetchList();
             }
        },
        methods: {
            fetchFocus() {
                fetchFocus().then(res => {
                    this.features = res.data || []
                }).catch(err => {
                    console.log(err)
                })
            },
            fetchList() {
                var ppp = {page:this.currPage+1}
                if(this.category_id !== undefined){
                    ppp["cate_id"] = this.category_id
                }
                if(this.searchWords){
                     ppp["title"] =this.searchWords
                }
                fetchList(ppp ).then(res => {
                    this.postList = this.postList.concat(res.data.list || [])
                    this.currPage = res.data.paginator.currpage
                    this.hasNextPage = res.data.paginator.currpage !== res.data.paginator.totalpages
                }).catch(err => {
                    console.log(err)
                })
                setting("IndexPics" ).then(res => {
              
                 }).catch(err => {
                    console.log(err)
                })
            },
            loadMore() {
                this.fetchList()
            }
        },
        created(){
           
        },
        mounted() {
                    this.features = []
                this.  postList = []
                this.  currPage =  0
                 this. hasNextPage = false
            this.fetchFocus();
            this.fetchList();
        }
    }
</script>
<style scoped lang="less">

    .site-content {
        .notify {
            margin: 60px 0;
            border-radius: 3px;
            & > div {
                padding: 20px;
            }
        }


        .search-result {
            padding: 15px 20px;
            text-align: center;
            font-size: 20px;
            font-weight: 400;
            border: 1px dashed #ddd;
            color: #828282;
        }
    }

    .top-feature {
        width: 100%;
        height: auto;
        margin-top: 30px;

        .feature-content {
            margin-top: 10px;
            display: flex;
            justify-content: space-between;
            position: relative;

            .feature-item {
                width: 32.9%;
            }
        }
    }

    .site-main {
        padding-top: 80px;

        &.search {
            padding-top: 0;
        }
    }

    .more{
        margin: 50px 0;
        .more-btn{
            width: 100px;
            height: 40px;
            line-height: 40px;
            text-align: center;
            color: #ADADAD;
            border: 1px solid #ADADAD;
            border-radius: 20px;
            margin: 0 auto;
            cursor: pointer;
            &:hover{
                color: #8fd0cc;
                border: 1px dashed #8fd0cc;
            }
        }
    }

    /******/
    @media (max-width: 800px) {
        .top-feature {
            display: none;
        }

        .site-main {
            padding-top: 40px;
        }

        .site-content {
            .notify {
                margin: 30px 0 0 0;
            }

            .search-result {
                margin-bottom: 20px;
                font-size: 16px;
            }
        }
    }

    /******/
</style>
