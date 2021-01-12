<template>
  <div class="createPost-container">

    <el-form ref="postForm" :model="postForm" :rules="rules" class="form-container">

      <sticky :z-index="10" :class-name="'sub-navbar '+postForm.status">
 
        <el-button v-loading="loading" style="margin-left: 10px;" type="success" @click="submitForm">
          Publish
        </el-button>
        <el-button v-loading="loading" type="warning" @click="draftForm">
          添加内容
        </el-button>
      </sticky>

      <div class="createPost-main-container">
        <el-row>
          <Warning />

          <el-col :span="24">
            <el-form-item style="margin-bottom: 40px;" prop="title">
              <MDinput v-model="postForm.title" :maxlength="100" name="name" required>
                标题
              </MDinput>
            </el-form-item>

            <div class="postInfo-container">
              <el-row>

                <el-col :span="10">
                  <el-form-item label-width="80px" label="文章分类:" class="postInfo-container-item">
                        <el-select v-model="postForm.cate_id" :placeholder="$t('table.cate_id')" clearable class="filter-item" required>
                          <el-option v-for="(item,index) in catelistData" :key="index" :label="item.name" :value="item.id" />
                        </el-select>
                  </el-form-item>
                </el-col>

                <el-col :span="10">
                  <el-form-item label-width="80px" label="状态:" class="postInfo-container-item">
                        <el-select v-model="postForm.status" :placeholder="$t('table.status')" clearable class="filter-item" required>
                          <el-option v-for="(item,index) in statusMap" :key="index" :label="item.name" :value="item.value" />
                        </el-select>
                  </el-form-item>
                </el-col>

                <el-col :span="5">
                     <el-form-item  label-width="50px" label="排序:">
                          <el-input v-model="postForm.sort" :rows="1" type="textarea" class="article-textarea" autosize placeholder="越大越靠前" />
                        </el-form-item>
                </el-col>



              </el-row>
            </div>
          </el-col>
        </el-row>

         <el-form-item style="margin-bottom: 40px;" label-width="70px" label="简短介绍:">
          <el-input v-model="postForm.summary"   type="textarea" class="article-textarea" autosize placeholder="Please enter the content" />
           <span v-show="contentShortLength" class="word-counter">{{ contentShortLength }}words</span>
        </el-form-item>

        <el-form-item prop="image" style="margin-bottom: 30px;">
          <Upload v-model="postForm.image" />
        </el-form-item>
 
      
             <el-form-item style="margin-bottom: 40px;"   v-for="item in postForm. content">
              <MDinput  v-model="item .title"  :maxlength="100" required>
                标签
              </MDinput>
                 
              <el-card style="height: 610px;margin-top:10px;">
                  <quill-editor v-model="item.content" ref="myQuillEditor" style="height: 500px;  "  :options="editorOption">
                  </quill-editor>
              </el-card>


            </el-form-item>
 

     
      </div>
    </el-form>
  </div>
</template>

<script>
import Tinymce from '@/components/Tinymce'
import Upload from '@/components/Upload/SingleImage3'
import MDinput from '@/components/MDinput'
import Sticky from '@/components/Sticky' // 粘性header组件
import { validURL } from '@/utils/validate'
import { fetchArticle, cateList,createArticle,updateArticle } from '@/api/article'

import { searchUser } from '@/api/remote-search'
import Warning from './Warning'
import { CommentDropdown,  SourceUrlDropdown } from './Dropdown'


  import {
    quillEditor
  } from 'vue-quill-editor'
  import 'quill/dist/quill.core.css'
  import 'quill/dist/quill.snow.css'
  import 'quill/dist/quill.bubble.css'


const defaultForm = {
  status: 0,
  title: '', // 文章题目
  content: [], // 文章内容
  id: undefined,
  cate_id:undefined,
  sort:0,
  summary:"",
  image:undefined,
}

const statusMap = [
  {value:1,name:'发布'},
  {value:0,name:'未发布'}
]



export default {
  name: 'ArticleDetail',
  components: { Tinymce, MDinput, Upload, Sticky, Warning, CommentDropdown,  SourceUrlDropdown, quillEditor },
  props: {
    isEdit: {
      type: Boolean,
      default: false
    }
  },
  data() {
    const validateRequire = (rule, value, callback) => {
      if (value === '') {
        this.$message({
          message: rule.field + '为必传项',
          type: 'error'
        })
        callback(new Error(rule.field + '为必传项'))
      } else {
        callback()
      }
    }
    const validateSourceUri = (rule, value, callback) => {
      if (value) {
        if (validURL(value)) {
          callback()
        } else {
          this.$message({
            message: '外链url填写不正确',
            type: 'error'
          })
          callback(new Error('外链url填写不正确'))
        }
      } else {
        callback()
      }
    }
    return {
      postForm: Object.assign({}, defaultForm),
      statusMap:statusMap,
      loading: false,
      userListOptions: [],
      contentList:[ ],
      rules: {
        image: [{ validator: validateRequire }],
        title: [{ validator: validateRequire }],
        content: [{ validator: validateRequire }],
        source_uri: [{ validator: validateSourceUri, trigger: 'blur' }]
      },
      tempRoute: {},
      catelistData:[],
        editorOption:  {}
         
 
    }
  },
  computed: {
    contentShortLength() {
      return this.postForm.summary.length
    },
    lang() {
      return this.$store.getters.language
    },
    displayTime: {
      // set and get is useful when the data
      // returned by the back end api is different from the front end
      // back end return => "2013-06-25 06:59:25"
      // front end need timestamp => 1372114765000
      get() {
        return (+new Date(this.postForm.display_time))
      },
      set(val) {
        this.postForm.display_time = new Date(val)
      }
    }
  },
  created() {
    if (this.isEdit) {
      const id = this.$route.params && this.$route.params.id
      this.fetchData(id)
    }

    cateList().then(response => {
          this.catelistData  =   response.data 
         
    })

    // Why need to make a copy of this.$route here?
    // Because if you enter this page and quickly switch tag, may be in the execution of the setTagsViewTitle function, this.$route is no longer pointing to the current page
    // https://github.com/PanJiaChen/vue-element-admin/issues/1221
    this.tempRoute = Object.assign({}, this.$route)
  },
  methods: {
    fetchData(id) {
      fetchArticle(id).then(response => {
        this.postForm = response.data
          if(this.postForm . content == null )  {
            this.postForm . content = []
          }
        
        // set tagsview title
        this.setTagsViewTitle()

        // set page title
        this.setPageTitle()
      }).catch(err => {
        console.log(err)
      })
    },
    setTagsViewTitle() {
      const title = this.lang === 'zh' ? '编辑文章' : 'Edit Article'
      const route = Object.assign({}, this.tempRoute, { title: `${title}-${this.postForm.id}` })
      this.$store.dispatch('tagsView/updateVisitedView', route)
    },
    setPageTitle() {
      const title = 'Edit Article'
      document.title = `${title} - ${this.postForm.id}`
    },
    submitForm() {
      this.$refs.postForm.validate(valid => {
        if (valid) {
          this.loading = true
          if (this.isEdit) {
            updateArticle(this.postForm).then(() => {
              this.loading = false
              this.$notify({
                title: 'Success',
                dangerouslyUseHTMLString: true,
                message: '修改成功',
                type: 'success',
                duration: 2000
              })
                    this.$router.push({path:'/articles/list' });
            })
          }else{
            createArticle(this.postForm).then(() => {
              this.loading = false
              this.$notify({
                title: 'Success',
                dangerouslyUseHTMLString: true,
                message: '添加成功',
                type: 'success',
                duration: 2000
              })
                this.$router.push({path:'/articles/list' });
            })


          }


          // this.$notify({
          //   title: '成功',
          //   message: '发布文章成功',
          //   type: 'success',
          //   duration: 2000
          // })
          // this.postForm.status = 'published'
          this.loading = false
        } else {
          console.log('error submit!!')
          return false
        }




      })
    },
    draftForm() {


    this. postForm. content.push({
        "title":"",
        "content":""
      })  
      return
      if (this.postForm.content.length === 0 || this.postForm.title.length === 0) {
        this.$message({
          message: '请填写必要的标题和内容',
          type: 'warning'
        })
        return
      }

     



      this.$message({
        message: '保存成功',
        type: 'success',
        showClose: true,
        duration: 1000
      })
      this.postForm.status =0
      // this.postForm.status = 'draft'
    },
    getRemoteUserList(query) {
      searchUser(query).then(response => {
        if (!response.data.items) return
        this.userListOptions = response.data.items.map(v => v.name)
      })
    }
  }
}
</script>

<style lang="scss" scoped>
@import "~@/styles/mixin.scss";

.createPost-container {
  position: relative;

  .createPost-main-container {
    padding: 40px 45px 20px 50px;

    .postInfo-container {
      position: relative;
      @include clearfix;
      margin-bottom: 10px;

      .postInfo-container-item {
        float: left;
      }
    }
  }

  .word-counter {
    width: 40px;
    position: absolute;
    right: 10px;
    top: 0px;
  }
}

.article-textarea /deep/ {
  textarea {
    padding-right: 40px;
    resize: none;
    border: none;
    border-radius: 0px;
    border-bottom: 1px solid #bfcbd9;
  }
}
</style>
