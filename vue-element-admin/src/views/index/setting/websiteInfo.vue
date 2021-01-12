<template>
  <div class="createPost-container">
    <el-form ref="postForm" :model="postForm" :rules="rules" class="form-container">
      <sticky :z-index="10" :class-name="'sub-navbar '+postForm.status">
        <el-button v-loading="loading" style="margin-left: 10px;" type="success" @click="submitForm">
          Publish
        </el-button>
      </sticky>

      <div class="createPost-main-container">
        <el-row>
          
          <el-col :span="24">
            <el-form-item style="margin-bottom: 40px;" prop="title">
              <MDinput v-model="postForm.slogan" :maxlength="60" name="slogan" required>
                slogan
              </MDinput>
            </el-form-item>

            <el-form-item style="margin-bottom: 40px;" prop="title">
              <MDinput v-model="postForm.name" :maxlength="20" name="name" required>
                name
              </MDinput>
            </el-form-item>  


            <el-form-item style="margin-bottom: 40px;" prop="title">
              <MDinput v-model="postForm.domain" :maxlength="50" name="domain" required>
                domain
              </MDinput>
            </el-form-item>  

            <el-form-item style="margin-bottom: 40px;" prop="title">
              <MDinput v-model="postForm.notice" :maxlength="50" name="notice" required>
                notice
              </MDinput>
            </el-form-item>  

            <el-form-item style="margin-bottom: 40px;" prop="title">
              <MDinput v-model="postForm.desc" :maxlength="30" name="desc" required>
                desc
              </MDinput>
            </el-form-item>  


          </el-col>
        </el-row>

     

        <el-form-item prop="image" style="margin-bottom: 30px;">
          <Upload v-model="postForm.avatar" />
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
import { PostWebsiteInfo , GetWebsiteInfo } from '@/api/setting'

import { searchUser } from '@/api/remote-search'
 
 


  import {
    quillEditor
  } from 'vue-quill-editor'
  import 'quill/dist/quill.core.css'
  import 'quill/dist/quill.snow.css'
  import 'quill/dist/quill.bubble.css'


const defaultForm = {
      avatar: 'https://s2.ax1x.com/2020/01/17/1SCadg.png',
      slogan: 'The way up is not crowded, and most chose ease.',
      name: 'FZY′blog',
      domain: 'https://www.fengziy.cn',
      notice: '本博客的Demo数据由Mockjs生成',
      desc: '一个It技术的探索者'
}

 


export default {
  name: 'Websiteinfo',
  components: { Tinymce, MDinput, Upload, Sticky,     quillEditor },
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
   
      loading: false,
  
      avatar: 'https://s2.ax1x.com/2020/01/17/1SCadg.png',
      slogan: 'The way up is not crowded, and most chose ease.',
      name: 'FZY′blog',
      domain: 'https://www.fengziy.cn',
      notice: '本博客的Demo数据由Mockjs生成',
      desc: '一个It技术的探索者',


  
      rules: {
        avatar: [{ validator: validateRequire }],
        slogan: [{ validator: validateRequire }],
        name: [{ validator: validateRequire }],
        domain: [{ validator: validateRequire }],
        notice: [{ validator: validateRequire }],
         desc: [{ validator: validateRequire }]
 
      },
       
       
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
     
  },
  created() {
    GetWebsiteInfo().then(response => {

  this. postForm =  response.data 
    })


  },
  methods: {
    fetchData(id) {
      fetchArticle(id).then(response => {
        this.postForm = response.data
        this.postForm.comment_disabled = response.data.comment_disabled==0?false:true
        this.postForm.author=response.data.authorname
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
      
            PostWebsiteInfo({websiteInfo:this.postForm}  ).then(() => {
              this.loading = false
              this.$notify({
                title: 'Success',
                dangerouslyUseHTMLString: true,
                message: '修改成功',
                type: 'success',
                duration: 2000
              })
               ///     this.$router.push({path:'/articles/list' });
            })
          



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
     
    },
    getRemoteUserList(query) {
      
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