<template>
  <div class="app-container">
    <div style="height:40px;float:right;margin-left:25px;margin-right:20px">
      <el-form :inline="true" label-width="80px">
        <el-form-item label="操作用户">
          <el-input
            size="mini"
            @keyup.enter.native="startgetoperatelog"
            v-model="operatequery.username"
          ></el-input>
        </el-form-item>
        <el-form-item label="操作类型">
          <el-select
            size="small"
            filterable
            clearable
            v-model="operatequery.method"
            @change="startgetoperatelog"
          >
            <el-option
              v-for="item in getoption(operatetype)"
              :key="item.label"
              :label="item.label"
              :value="item.value"
            ></el-option>
          </el-select>
        </el-form-item>
        <el-form-item label="模块">
          <el-select
            size="small"
            filterable
            clearable
            v-model="operatequery.module"
            @change="startgetoperatelog"
          >
            <el-option
              v-for="item in getoption(moduletype)"
              :key="item.label"
              :label="item.label"
              :value="item.value"
            ></el-option>
          </el-select>
        </el-form-item>
      </el-form>
    </div>
    <el-table
      v-loading="listLoading"
      :data="data"
      stripe
      fit
      highlight-current-row
      style="width: 100%;"
    >
      <el-table-column align="center" fixed="left" label="操作用户" min-width="80">
        <template slot-scope="scope">
          <span>{{ scope.row.user_name }}</span>
        </template>
      </el-table-column>

      <el-table-column align="center" label="用户类型" min-width="100">
        <template slot-scope="scope">
          <span v-if="scope.row.user_role === 1">普通用户</span>
          <span v-else-if="scope.row.user_role === 2">管理员</span>
          <span v-else-if="scope.row.user_role === 3">访客</span>
          <span v-else>Unknow</span>
        </template>
      </el-table-column>
      <el-table-column align="center" label="操作类型" min-width="100">
        <template slot-scope="scope">
          <span v-if="scope.row.method === 'PUT'">
            <el-tag type="warning" size="small">{{ operatetype[scope.row.method] }}</el-tag>
          </span>
          <span v-else-if="scope.row.method === 'DELETE'">
            <el-tag type="danger" size="small">{{ operatetype[scope.row.method] }}</el-tag>
          </span>
          <span v-else-if="scope.row.method === 'POST'">
            <el-tag type="success" size="small">{{ operatetype[scope.row.method] }}</el-tag>
          </span>
          <span v-else>
            <el-tag type="info" size="small">Unknow</el-tag>
          </span>
        </template>
      </el-table-column>

      <el-table-column align="center" label="修改模块" min-width="70">
        <template slot-scope="scope">
          <span>{{ moduletype[scope.row.module] }}</span>
        </template>
      </el-table-column>
      <el-table-column align="center" label="修改模块名称" min-width="100">
        <template slot-scope="scope">
          <span>{{ scope.row.module_name }}</span>
        </template>
      </el-table-column>
      <el-table-column align="center" label="操作时间" min-width="150">
        <template slot-scope="scope">
          <span>{{ scope.row.operate_time }}</span>
        </template>
      </el-table-column>
      <el-table-column align="center" label="描述" min-width="120">
        <template slot-scope="scope">
          <span v-if="scope.row.desc === ''">-</span>
          <span v-else>{{ scope.row.desc }}</span>
        </template>
      </el-table-column>
      <el-table-column align="center" label="操作" min-width="100">
        <template slot-scope="scope">
          <el-button
            :disabled="scope.row.columns.length === 0"
            type="primary"
            size="mini"
            @click="lookdetail(scope.row.columns)"
          >详情</el-button>
        </template>
      </el-table-column>
    </el-table>
    <div style="margin-top: 10px;float:right;height: 70px;">
      <el-pagination
        :page-size="operatequery.limit"
        @current-change="handleCurrentChangerun"
        background
        layout="total,prev, pager, next"
        :total="pagecount"
      ></el-pagination>
    </div>

    <el-dialog title="操作详情" :visible.sync="diaogVisible" center width="60%">
      <el-table
        v-loading="listLoading"
        :data="detaildata"
        stripe
        size="mini"
        border
        fit
        highlight-current-row
        style="width: 100%;"
      >
        <el-table-column align="center" fixed="left" label="操作字段" min-width="70">
          <template slot-scope="scope">
            <span>{{ scope.row.name }}</span>
          </template>
        </el-table-column>
        <el-table-column align="center" fixed="left" label="旧值" min-width="100">
          <template slot-scope="scope">
            <span v-if="scope.row.old_value === null">-</span>
            <span v-else>{{ scope.row.old_value }}</span>
          </template>
        </el-table-column>
        <el-table-column align="center" fixed="left" label="新值" min-width="100">
          <template slot-scope="scope">
            <span v-if="scope.row.new_value === null">-</span>
            <span v-else v-text="scope.row.new_value"></span>
          </template>
        </el-table-column>
      </el-table>
    </el-dialog>
  </div>
</template>

<script>
import { Message } from "element-ui";
import router from "@/router";

import { getoperatelog } from "@/api/user";

export default {
  data() {
    return {
      diatasktitle: "",
      diaogVisible: false,
      tasklogtitle: "",
      tasklog: "",
      taskselect: [],
      data: [],
      pagecount: 0,
      detaildata: [],
      listLoading: false,
      operatequery: {
        username: "",
        method: "",
        module: "",
        offset: 0,
        limit: 15
      },
      operatetype: {
        PUT: "修改",
        POST: "创建",
        DELETE: "删除"
      },
      moduletype: {
        user: "用户",
        task: "任务",
        hostgroup: "主机组",
        host: "主机"
      },
      defaultProps: {
        children: "children",
        label: "name"
      }
    };
  },

  created() {
    this.startgetoperatelog();
  },
  methods: {
    getoption(datatype) {
      var dataoption = [];
      Object.entries(datatype).map(function(key) {
        dataoption.push({
          value: key[0],
          label: key[1]
        });
      });

      return dataoption;
    },
    handleCurrentChangerun(page) {
      this.operatequery.offset = (page - 1) * this.operatequery.limit;
      this.startgetoperatelog();
    },
    startgetoperatelog() {
      getoperatelog(this.operatequery).then(response => {
        this.data = response.data;
        this.pagecount = response.count;
      });
    },
    lookdetail(data) {
      this.diaogVisible = true;
      this.detaildata = data;
    }
  }
};
</script>


<style lang="scss" scoped>
.el-button--mini,
.el-button--mini.is-round {
  padding: 5px 5px;
}

.demo-table-expand {
  font-size: 0;
}
.demo-table-expand label {
  width: 90px;
  color: #99a9bf;
}
.demo-table-expand .el-form-item {
  margin-right: 0;
  margin-bottom: 0;
  width: 50%;
}
.el-button--mini,
.el-button--mini.is-round {
  padding: 5px 5px;
}
.sub-title {
  text-align: left;
  color: #909399;
  font-size: 20px;
  font-family: "Helvetica Neue, Helvetica, PingFang SC, Hiragino Sans GB, Microsoft YaHei, Arial, sans-serif";
  // margin-bottom: 6px;
  font-weight: 700;
}
.string {
  color: green;
}
.number {
  color: darkorange;
}
.boolean {
  color: blue;
}
.null {
  color: magenta;
}
.key {
  color: red;
}
</style>


