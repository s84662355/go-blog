<template>
  <div class="app-container">
    <div class="filter-container">

      <el-input v-model="listQuery.title" :placeholder="$t('table.keyword')" style="width: 200px;" class="filter-item" @keyup.enter.native="handleFilter" />


      <el-select v-model="listQuery.cate_id" :placeholder="$t('table.cate_id')" clearable class="filter-item" style="width: 130px">
        <el-option v-for="(item,index) in catelistData" :key="index" :label="item.name" :value="item.id" />
      </el-select>


      <el-select v-model="listQuery.status" :placeholder="$t('table.status')" clearable class="filter-item" style="width: 130px">
        <el-option v-for="(item,index) in statusMap" :key="index" :label="item.name" :value="item.value" />
      </el-select>


       
      <el-date-picker
        v-model="listQuery.dateValue"
        class="filter-item"
        type="datetimerange"
        start-placeholder="开始日期"
        end-placeholder="结束日期"
        value-format="yyyy-MM-dd HH:mm:ss"
        :default-time="['00:00:00','23:59:59']"
      />
      <el-button v-waves class="filter-item" type="primary" icon="el-icon-search" @click="handleFilter">
        {{ $t('table.search') }}
      </el-button>
      <el-button v-waves :loading="downloadLoading" class="filter-item" type="primary" icon="el-icon-download" @click="handleDownload">
        {{ $t('table.export') }}
      </el-button>
    </div>
    <el-table v-loading="listLoading" :data="list" border fit highlight-current-row  >
      <el-table-column align="center" label="id" width="60px"   >
        <template slot-scope="scope">
          <span>{{ scope.row.id }}</span>
        </template>
      </el-table-column>


      <el-table-column min-width="100px"  align="center" label="标题">
        <template slot-scope="scope">
          <span>{{ scope.row.title }}</span>
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

 

      <el-table-column width="100px" align="center" label="文章分类">
        <template slot-scope="scope">
          <span>{{ scope.row.cate_name }}</span>
        </template>
      </el-table-column>


      <el-table-column width="50px" align="center" label="文章阅读次数">
        <template slot-scope="scope">
          <span>{{ scope.row.read_amount }}</span>
        </template>
      </el-table-column>


      <el-table-column width="80px" align="center" label="文章排序">
        <template slot-scope="scope">
          <span>{{ scope.row.sort}}</span>
        </template>
      </el-table-column>

      <el-table-column   class-name="status-col" label="状态" width="80px"   >
        <template slot-scope="{row}">
          <el-tag :type="row.status | statusFilter">
            {{ row.status == 0 ? "未发布" : "发布" }}
          </el-tag>
        </template>
      </el-table-column>

 

      <el-table-column min-width="150px"  align="center" label="操作"  >
        <template slot-scope="scope">
            <router-link :to="'/articles/edit/'+scope.row.id">
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

    <pagination v-show="total>0" :total="total" :page.sync="listQuery.page" :limit.sync="listQuery.limit" @pagination="getList" />
  </div>
</template>

<script>
import {articleList, cateList,fetchList } from '@/api/article'
import waves from '@/directive/waves' // waves directive
import { parseTime } from '@/utils'
import Pagination from '@/components/Pagination' // Secondary package based on el-pagination

const statusMap = [
  {value:1,name:'发布'},
  {value:0,name:'未发布'}
]
 
export default {
  name: 'ArticleList',
  directives: { waves },
  components: { Pagination },
  filters: {
    statusFilter(status) {
      const statusMap = {
        1: 'success',
        0: 'info',
        10: 'danger'
      }
      return statusMap[status]
    },
    statusNameFilter(status) {
      return statusMap[status]
    }
  },
  data() {
    return {
      list: null,
      total: 0,
      listLoading: true,
      listQuery: {
        page: 1,
        limit: 10,
        importance: undefined,
        title: undefined,
        status: undefined,
        dateValue: '',
        cate_id:''
      },
      importanceOptions: [1, 2, 3],
      statusMap,
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
    this.getList()
 
    cateList().then(response => {
          this.catelistData  =   response.data 
         
    })
  },
  methods: {
    handleFilter() {
      this.listQuery.page = 1
      this.getList()
    },
    resetTemp() {
      this.temp = {
        id: undefined,
        importance: 1,
        remark: '',
        timestamp: new Date(),
        title: '',
        status: 'published',
        type: ''
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
    },
    handleDownload() {
      this.downloadLoading = true
      import('@/vendor/Export2Excel').then(excel => {
        const tHeader = ['display_time', 'title', 'importance', 'status']
        const filterVal = ['display_time', 'title', 'importance', 'status']
        const data = this.formatJson(filterVal, this.list)
        excel.export_json_to_excel({
          header: tHeader,
          data,
          filename: 'table-list'
        })
        this.downloadLoading = false
      })
    },
    getList() {
      this.listLoading = true


/*
      fetchList(this.listQuery).then(response => {
         this.list = response.data.items
         this.total = response.data.total
         this.listLoading = false
      })
*/

 
 

      articleList(this.listQuery).then(response => {
         this.list = response.data.list
        this.listLoading = false
         this.total = response.data.paginator.total
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
