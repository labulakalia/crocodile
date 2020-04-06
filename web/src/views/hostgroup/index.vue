<template>
  <div class="app-container">
    <div v-if="is_change || is_create" style="margin-left:25px;margin-right:80px；height:80px">
      <el-form
        :model="hostgroup"
        ref="hostgroup"
        :rules="rules"
        label-position="right"
        label-width="120px"
        size="mini"
      >
        <el-form-item label="执行器名称" prop="name">
          <el-input
            :disabled="is_change"
            v-model="hostgroup.name"
            placeholder="请输入执行器名称"
            clearable
            style="width: 500px;"
          ></el-input>
        </el-form-item>
        <el-form-item label="Worker" prop="task_type" placeholder="请选择主机">
          <el-select multiple filterable v-model="hostgroup.addrs" style="width: 500px;">
            <el-option
              v-for="item in hostselect"
              :key="item.label"
              :label="item.label"
              :value="item.value"
            >
              <span style="float: left">{{ item.label }}</span>
              <!-- <span
                style="float: right; color: #8492a6; font-size: 13px;margin-right: 30px;"
              >{{ item.online }}</span>-->
              <span style="float: right; color: #8492a6; font-size: 13px;margin-right: 30px;">
                <el-tag v-if="item.online === 1" size="mini" type="success">Online</el-tag>
                <el-tag v-if="item.online === -1" size="mini" type="danger">Offline</el-tag>
              </span>
            </el-option>
          </el-select>
        </el-form-item>
        <el-form-item label="备注" prop="remark">
          <el-input
            type="textarea"
            v-model="hostgroup.remark"
            placeholder="请输入任务备注"
            clearable
            style="width: 500px;"
          ></el-input>
        </el-form-item>
      </el-form>
      <div style="margin-left: 120px;">
        <el-button size="small" type="primary" @click="submithostgroup('hostgroup')">确 定</el-button>
        <el-button size="small" @click="is_create = false;is_change = false">取 消</el-button>
      </div>
    </div>
    <div v-else>
      <div style="float: right">
        <el-tooltip class="item" effect="dark" content="新建主机组" placement="top-start">
          <el-button type="primary" size="small" @click="createhostgrouppre">New</el-button>
        </el-tooltip>
      </div>
      <el-table
        v-loading="listLoading"
        :data="data"
        stripe
        fit
        highlight-current-row
        style="width: 100%;"
      >
        <el-table-column align="center" fixed="left" label="名称" min-width="100">
          <template slot-scope="scope">
            <span>{{ scope.row.name }}</span>
          </template>
        </el-table-column>

        <el-table-column align="center" label="主机" min-width="70">
          <template slot-scope="scope">
            <el-popover
              placement="right"
              width="700"
              trigger="click"
              @show="gethostdetail(scope.row)"
            >
              <el-table border :data="hghosts">
                <el-table-column property="addr" label="IP" width="150"></el-table-column>
                <el-table-column property="hostname" label="主机名" min-width="100"></el-table-column>
                <el-table-column property="weight" label="权重" width="60"></el-table-column>
                <el-table-column property="version" label="版本" width="60"></el-table-column>
                <el-table-column label="状态" width="70">
                  <template slot-scope="scope">
                    <el-tag type="success" size="mini" v-if="scope.row.online">Online</el-tag>
                    <el-tag type="danger" size="mini" v-else>Offline</el-tag>
                  </template>
                </el-table-column>
                <el-table-column label="暂停" width="80">
                  <template slot-scope="scope">
                    <el-tag v-if="!scope.row.stop" size="mini" type="success">Normal</el-tag>
                    <el-tag v-else size="mini" type="danger">Stop</el-tag>
                  </template>
                </el-table-column>
              </el-table>

              <el-button
                :disabled="scope.row.addrs.length === 0"
                type="success"
                size="small"
                slot="reference"
                round
              >
                <strong>{{ scope.row.addrs.length }}</strong>
              </el-button>
            </el-popover>
          </template>
        </el-table-column>

        <el-table-column align="center" label="创建人" min-width="70">
          <template slot-scope="scope">
            <span>{{ scope.row.create_by }}</span>
          </template>
        </el-table-column>
        <el-table-column property="create_time" label="创建时间" width="160"></el-table-column>
        <el-table-column property="update_time" label="修改时间" width="160"></el-table-column>
        <el-table-column align="center" label="备注" min-width="100">
          <template slot-scope="scope">
            <span>{{ scope.row.remark }}</span>
          </template>
        </el-table-column>
        <el-table-column fixed="right" align="center" label="操作" min-width="80">
          <template slot-scope="scope">
            <el-button-group>
              <el-button type="warning" size="mini" @click="changehostgrouppre(scope.row)">修改</el-button>
              <el-popconfirm
                :hideIcon="true"
                title="确定删除主机组?"
                @onConfirm="deletehostgrouppre(scope.row)"
              >
                <el-button slot="reference" type="danger" size="mini">删除</el-button>
              </el-popconfirm>
              <!-- <el-button type="danger" size="mini" @click="">删除</el-button> -->
            </el-button-group>
          </template>
        </el-table-column>
      </el-table>
      <div style="margin-top: 10px;float:right;height: 70px;">
        <el-pagination
          :page-size="hostgroupquery.limit"
          @current-change="handleCurrentChangerun"
          background
          layout="total,prev, pager, next"
          :total="pagecount"
        ></el-pagination>
      </div>
    </div>
  </div>
</template>

<script>
import {
  gethostgroup,
  getselecthostgroup,
  createhostgroup,
  deletehostgroup,
  changehostgroup,
  gethostsbyhgid
} from "@/api/hostgroup";

import { getselecthost } from "@/api/host";

import { Message } from "element-ui";

export default {
  data() {
    return {
      listLoading: false,
      data: [],
      pagecount: 0,
      executeips: [],
      createactuatordialog: false,
      updateactuatordialog: false,
      rules: {
        name: [{ required: true, message: "请输入执行器名称", trigger: "blur" }]
      },
      addips: [],
      arronline: [],
      hostgroupquery: {
        offset: 0,
        limit: 15
      },
      hostgroup: {
        id: "",
        name: "",
        addrs: [],
        remark: ""
      },
      hghosts: [],
      is_change: false,
      is_create: false,
      hostselect: []
    };
  },
  created() {
    this.getallhostgroup();
  },
  methods: {
    getallhostgroup() {
      gethostgroup(this.hostgroupquery).then(resp => {
        this.data = resp.data;
        this.pagecount = resp.count;
      });
    },
    createhostgrouppre() {
      this.startgetselecthost();
      if (this.hostgroup.name === undefined) {
        this.hostgroup["name"] = "";
      } else {
        this.hostgroup.name = "";
      }

      this.hostgroup.addrs = "";
      this.hostgroup.remark = "";
      this.is_create = true;
    },
    gethostdetail(hostgroup) {
      var query = {
        id: hostgroup.id
      };
      gethostsbyhgid(query).then(resp => {
        this.hghosts = resp.data;
      });
    },
    changehostgrouppre(hostgroup) {
      this.startgetselecthost();
      if (this.hostgroup.id === undefined) {
        this.hostgroup["id"] = hostgroup.id;
      } else {
        this.hostgroup.id = hostgroup.id;
      }
      this.hostgroup.name = hostgroup.name;

      this.hostgroup.addrs = hostgroup.addrs;
      this.hostgroup.remark = hostgroup.remark;
      this.is_change = true;
    },
    deletehostgrouppre(hostgroup) {
      var delid = {
        id: hostgroup.id
      };
      deletehostgroup(delid).then(resp => {
        if (resp.code === 0) {
          Message.success(`删除主机组 ${hostgroup.name} 成功`);
          this.getallhostgroup();
          this.is_create = false;
        } else {
          Message.error(`删除主机组 ${hostgroup.name} 失败: ${resp.msg}`);
        }
      });
    },
    startgetselecthost() {
      getselecthost().then(resp => {
        this.hostselect = resp.data;
      });
    },
    submithostgroup(formName) {
      this.$refs[formName].validate(valid => {
        if (valid) {
          var name = this.hostgroup.name;
          if (this.is_create === true) {
            delete this.hostgroup.id;
            createhostgroup(this.hostgroup).then(resp => {
              if (resp.code === 0) {
                Message.success(`创建主机组 ${name} 成功`);
                this.getallhostgroup();
                this.is_create = false;
              } else {
                Message.error(`创建主机组 ${name} 失败: ${resp.msg}`);
              }
            });
          }
          if (this.is_change === true) {
            delete this.hostgroup.name;
            changehostgroup(this.hostgroup).then(resp => {
              if (resp.code === 0) {
                Message.success(`修改主机组 ${name} 成功`);
                this.getallhostgroup();
                this.is_change = false;
              } else {
                Message.error(`修改主机组 ${name} 失败: ${resp.msg}`);
              }
            });
          }
        }
      });
    },
    handleCurrentChangerun(page) {
      this.hostgroupquery.offset = (page - 1) * this.hostgroupquery.limit;
      this.getallhostgroup();
    }
  }
};
</script>

<style scoped>
.el-input--suffix .el-input__inner {
  padding-right: 40px;
}
.el-button--mini,
.el-button--mini.is-round {
  padding: 5px 5px;
}
.card-panel-text {
  line-height: 18px;
  color: rgba(0, 0, 0, 0.45);
  font-size: 16px;
  margin-bottom: 12px;
}
</style>
