<template>
	<div style="float:left;margin-left:100px;">
  <el-form :model="form">
  	            
    <el-form-item label="图片" :label-width="formLabelWidth">
      <el-upload
        ref="upload"
        action="#"
        accept="image/png,image/gif,image/jpg,image/jpeg"
        list-type="picture-card"
        :limit=limitNum
        :auto-upload="true"
        :on-exceed="handleExceed"
        :before-upload="handleBeforeUpload"
        :on-preview="handlePictureCardPreview"
        :http-request="uploadSectionFile" 
        :file-list="fileList"
        :on-remove="handleRemove">
        <i class="el-icon-plus"></i>
      </el-upload>
      <el-dialog :visible.sync="dialogVisible">
         <img width="100%" :src="dialogImageUrl" alt="">
      </el-dialog>
    </el-form-item>
    <el-form-item>
 
          <el-button size="small" type="primary" @click="submitFile">提交</el-button>
      <el-button size="small" @click="clearFile">取消</el-button>
    </el-form-item>
  </el-form>
</div>
</template>

<script>

import axios from 'axios'
import { MessageBox, Message } from 'element-ui'
import store from '@/store'
import { getToken } from '@/utils/auth'
import { httphost } from '@/utils/global'
 
import { UpdateIndexPics,GetIndexPics} from '@/api/setting'

export default {
  name: 'SettingUplaod',
  data() {
    return{
      dialogImageUrl: '',
      dialogVisible: false,
      formLabelWidth: '80px',
      limitNum: 5,
      form: {},
      upload_url:httphost + '/upload/image?token='+getToken(),
      fileList : []
    }
  },
 
    computed: {
 
    },


  created() {


        GetIndexPics( ).then(response => {
          this.fileList  =  response.data 
 

        }).catch(err => {
          console.log(err)
        })
  },


  methods: {
 
  	clearFile(){
  		  this.fileList = []
  	},
  	submitFile(){
  		console.log(  this.fileList )
 
        UpdateIndexPics({imgs:this.fileList}).then(response => {
          this.postForm = response.data

                   this.$notify({
                title: 'Success',
                dangerouslyUseHTMLString: true,
                message: '修改成功',
                type: 'success',
                duration: 2000
              })
        }).catch(err => {
          console.log(err)
        })

  	},
  	uploadSectionFile (params) {
  		console.log(params)

  		    var self = this,
            file = params.file,
            fileType = file.type,         
            file_url = self.$refs.upload.uploadFiles[0].url;

            var img = new Image();
            img.src = file_url;
            img.onload = function () {               
              self.uploadFile11(file  );
            }
 

  	},

    uploadFile11: function (file ) {
        var self = this,
            formData = new FormData();
        formData.append( 'file', file);

 
        axios.post(self.upload_url, formData, { headers: { 'Content-Type': 'multipart/form-data' } })
            .then(function (res) {

            	if(res.data.code == 20000){
                     
 

                      self.fileList .push({
                      	name:res.data.data,
                      	url:res.data.data
                      })


            	}else{
            		  _.$alert('上传失败，请重新上传', '提示', { type: 'error' });
            	}

            	 
			    //  self.$refs.upload.uploadFiles = []; 
            })
            .catch(function (err) {
               
            });
         
    },



    // 上传文件之前的钩子
    handleBeforeUpload(file){
      console.log('before')
      if(!(file.type === 'image/png' || file.type === 'image/gif' || file.type === 'image/jpg' || file.type === 'image/jpeg')) {
        this.$notify.warning({
          title: '警告',
          message: '请上传格式为image/png, image/gif, image/jpg, image/jpeg的图片'
        })
      }
      let size = file.size / 1024 / 1024 / 2
      if(size > 2) {
        this.$notify.warning({
          title: '警告',
          message: '图片大小必须小于2M'
        })
      }
    },
    // 文件超出个数限制时的钩子
    handleExceed(files, fileList) {

    },
    // 文件列表移除文件时的钩子
    handleRemove(file, fileList) {
    	 this.fileList = fileList

      console.log( this.fileList);
    },
    // 点击文件列表中已上传的文件时的钩子
    handlePictureCardPreview(file) {
      this.dialogImageUrl = file.url;
      this.dialogVisible = true;
    },
   
    uploadFile() {
      this.$refs.upload.submit()
    }
     
  } 
}
</script>

<style lang="scss" scoped>

</style>

 

