<template>
  <div class="app-container">
    <div class="filter-container">
     
          <el-button v-waves class="filter-item" type="primary" icon="el-icon-create" @click="handleFilter">
            创建分类
          </el-button>
      
    
    </div>
    <el-table v-loading="listLoading" :data="catelistData" border fit highlight-current-row  >
      <el-table-column align="center" label="id" width="60px"   >
        <template slot-scope="scope">
          <span>{{ scope.row.id }}</span>
        </template>
      </el-table-column>


      <el-table-column min-width="100px"  align="center" label="分类">
        <template slot-scope="scope">
          <span>{{ scope.row.name }}</span>
        </template>
      </el-table-column>

      <el-table-column width="120px" align="center" label="创建日期">
        <template slot-scope="scope">
          <span>{{ scope.row.created_at   }}</span>
        </template>
      </el-table-column>
      <el-table-column width="120px" align="center" label="更新日期">
        <template slot-scope="scope">
          <span>{{ scope.row.updated_at }}</span>
        </template>
      </el-table-column>

  

      <el-table-column min-width="150px"  align="center" label="操作"  >
        <template slot-scope="scope">
            <router-link :to="'/cate/edit/'+scope.row.id">
              <el-button type="primary" size="small" icon="el-icon-edit">
                编辑
              </el-button>
            </router-link>
            
              <el-button type="primary" size="small" icon="el-icon-delete" @click="delete(scope.row.id)">
                  删除
              </el-button>
         
        </template>
      </el-table-column>


    </el-table>
 
  </div>
</template>

<script>
import {articleList, cateList,fetchList } from '@/api/article'
import waves from '@/directive/waves' // waves directive
import { parseTime } from '@/utils'
import Pagination from '@/components/Pagination' // Secondary package based on el-pagination

 
export default {
  name: 'CateList',
  directives: { waves },
  components: { Pagination },
  filters: {
   
  },
  data() {
    return {
      list: null,
      total: 0,
      listLoading: true,
      listQuery: {
        page: 1,
        limit: 10,
        name: undefined,     
      
      },
      
    
      dialogFormVisible: false,
      dialogStatus: '',
      textMap: {
        update: 'Edit',
        create: 'Create'
      },
      dialogPvVisible: false,
      downloadLoading: false,
       catelistData:[],
    }
  },
  created() {
    
     this.listLoading = true
    cateList().then(response => {
          this.catelistData  =   response.data 
          this.listLoading = false
         
    })
  },
  methods: {
    handleFilter() {
      this.$router.push({path:'/cate/create' })
    },
    resetTemp() {
      this.temp = {
        id: undefined,
        timestamp: new Date(),
        title: '',
        
      }
    },
    formatJson(filterVal, jsonData) {
      return jsonData.map(v => filterVal.map(j => {
        if (j === 'timestamp') {
          return parseTime(v[j])
        } else {
          return v[j]
        }
      }))
    },
    handleCreate() {
      this.resetTemp()
      this.dialogStatus = 'create'
      this.dialogFormVisible = true
      this.$nextTick(() => {
        this.$refs['dataForm'].clearValidate()
      })
    } 
  }
}
</script>

<style scoped>
.edit-input {
  padding-right: 100px;
}
.cancel-btn {
  position: absolute;
  right: 15px;
  top: 10px;
}
</style>
