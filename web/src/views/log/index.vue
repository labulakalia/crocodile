<template>
  <div class="app-container">
    <div style="height:40px;float:right;margin-left:25px;margin-right:20px">
      <el-form :inline="true" label-width="80px">
        <el-form-item label="任务名称">
          <el-input
            size="mini"
            @keyup.enter.native="startgettasklog"
            clearable
            v-model="logquery.name"
          ></el-input>
        </el-form-item>
        <el-form-item label="执行结果">
          <el-select size="mini" v-model="logquery.status" @change="startgettasklog">
            <el-option
              v-for="item in statusoption"
              :key="item.label"
              :label="item.label"
              :value="item.value"
            ></el-option>
          </el-select>
        </el-form-item>
        <el-form-item>
          <el-tooltip content="清理运行任务日志" effect="dark" placement="top-start">
            <el-button type="warning" size="mini" @click="cleanlogvisible=true">清理</el-button>
          </el-tooltip>
        </el-form-item>
      </el-form>
      <el-dialog title="清理日志" center :visible.sync="cleanlogvisible" width="25%">
        <el-form>
          <el-form-item label="任务名称">
            <el-input
              v-model="logquery.name"
              size="mini"
              placeholder="请输入需要清理的任务名称"
              style="width: 70%;"
            ></el-input>
          </el-form-item>
          <el-form-item label="清理时间">
            <el-select size="mini" v-model="cleanpredat" style="width: 70%;">
              <el-option
                v-for="item in cleantimeoption"
                :key="item.label"
                :label="item.label"
                :value="item.value"
              ></el-option>
            </el-select>
          </el-form-item>
        </el-form>

        <p></p>
        <div style="text-align: center;">
          <el-button type="primary" size="mini" @click="startcleantask">确定</el-button>
          <el-button size="mini" type="text" @click="cleanlogvisible = false">取消</el-button>
        </div>
      </el-dialog>
    </div>
    <el-table
      v-loading="listLoading"
      :data="data"
      stripe
      fit
      highlight-current-row
      style="width: 100%;"
    >
      <el-table-column align="center" fixed="left" label="任务名称" min-width="80">
        <template slot-scope="scope">
          <span>{{ scope.row.name }}</span>
        </template>
      </el-table-column>

      <el-table-column align="center" label="开始时间" min-width="155">
        <template slot-scope="scope">
          <span>{{ scope.row.start_timestr }}</span>
        </template>
      </el-table-column>
      <el-table-column align="center" label="结束时间" min-width="155">
        <template slot-scope="scope">
          <span>{{ scope.row.end_timestr }}</span>
        </template>
      </el-table-column>

      <el-table-column align="center" label="运行时间" min-width="100">
        <template slot-scope="scope">
          <span v-if="scope.row.total_runtime <= 1000">{{ scope.row.total_runtime }}ms</span>
          <span v-if="scope.row.total_runtime > 1000">{{ scope.row.total_runtime/1000 }}s</span>
        </template>
      </el-table-column>
      <el-table-column align="center" label="触发方式" min-width="100">
        <template slot-scope="scope">
          <span>{{ scope.row.trigger_str }}</span>
        </template>
      </el-table-column>
      <el-table-column align="center" label="执行结果" min-width="80">
        <template slot-scope="scope">
          <el-tag v-if="scope.row.status === 1" type="success">成功</el-tag>
          <el-tag v-else-if="scope.row.status === -1" type="danger">失败</el-tag>
        </template>
      </el-table-column>
      <el-table-column align="center" label="操作" width="150">
        <template slot-scope="scope">
          <el-button-group>
            <el-button
              v-if="scope.row.status === -1"
              type="danger"
              size="mini"
              @click="gettasklogpre(scope.row)"
            >详情</el-button>
            <el-button v-else type="success" size="mini" @click="gettasklogpre(scope.row)">详情</el-button>
          </el-button-group>
        </template>
      </el-table-column>
      <el-table-column fixed="right" type="expand">
        <template slot-scope="scope">
          <div v-if="scope.row.status === -1">
            <span>
              <strong>ErrTask:</strong>
              {{ scope.row.err_task }}
            </span>
            <!-- <br /> -->
            <el-divider></el-divider>
            <span>
              <strong>ErrCode:</strong>
              {{ scope.row.err_code }}
            </span>
            <!-- <br /> -->
            <el-divider></el-divider>
            <span>
              <strong>ErrMsg:</strong>
              {{ scope.row.err_msg }}
            </span>

            <el-divider></el-divider>
            <span>
              <strong>ErrTaskType:</strong>
              {{ tasktype[scope.row.err_tasktype] }}
            </span>
          </div>
        </template>
      </el-table-column>
    </el-table>
    <div style="margin-top: 10px;float:right;height: 70px;">
      <el-pagination
       :page-size="logquery.limit"
        @current-change="handleCurrentChangerun"
        background
        layout="total,prev, pager, next"
        :total="pagecount"
      ></el-pagination>
    </div>

    <el-dialog :title="diatasktitle" :visible.sync="diaogVisible" center width="70%">
      <el-container>
        <el-aside width="200px">
          <el-tree
            :data="treetaskdata"
            highlight-current
            :default-expand-all="true"
            :props="defaultProps"
            :render-content="renderContent"
            @node-click="handleNodeClick"
          ></el-tree>
        </el-aside>
        <el-main>
          <span class="sub-title" v-html="tasklogtitle"></span>
          <el-card v-if="tasklog != ''" :body-style="{ padding: '0px' }" style="margin-top: 5px;">
            <editor
              v-model="tasklog"
              theme="solarized_dark"
              lang="text"
              height="500"
              width="100%"
              @init="initEditor"
              :options="{ readOnly: true }"
            ></editor>
          </el-card>
        </el-main>
      </el-container>
    </el-dialog>
  </div>
</template>

<script>
import { gettaskLog, gettaskLogTree, cleantasklog } from "@/api/task";
import { Message } from "element-ui";
import router from "@/router";

export default {
  components: {
    editor: require("vue2-ace-editor")
  },
  data() {
    return {
      diatasktitle: "",
      diaogVisible: false,
      tasklogtitle: "",
      tasklog: "",
      taskselect: [],
      data: [],
      treetaskdata: [],
      pagecount: 0,
      listLoading: false,
      logquery: {
        name: "",
        offset: 0,
        limit: 15,
        status: 0
      },
      statusoption: [
        {
          label: "成功",
          value: 1
        },
        {
          label: "失败",
          value: -1
        },
        {
          label: "全部",
          value: 0
        }
      ],
      tasktype: {
        1: "主任务",
        2: "父任务",
        3: "子任务"
      },
      defaultProps: {
        children: "children",
        label: "name"
      },
      cleanlogvisible: false,
      cleanpredat: "",
      cleantimeoption: [
        {
          label: "全部日志",
          value: 0
        },
        {
          label: "一周以前的日志",
          value: 7
        },

        {
          label: "一月以前的日志",
          value: 30
        },
        {
          label: "两个月以前的日志",
          value: 60
        },
        {
          label: "三个月以前的日志",
          value: 90
        },
        {
          label: "半年以前的日志",
          value: 182
        },
        {
          label: "一年以前的日志",
          value: 365
        }
      ]
    };
  },

  created() {
    // 刚进入页面是查询当前到一天前的日志
    if (this.$route.query.name !== undefined) {
      this.logquery.name = this.$route.query.name;
    }

    this.startgettasklog();
  },
  methods: {
    initEditor: function(editor) {
      require("brace/ext/language_tools");
      require("brace/mode/text");
      require("brace/theme/solarized_dark");
    },

    handleCurrentChangerun(page) {
      this.logquery.offset = (page - 1) * this.logquery.limit;
      this.startgettasklog();
    },
    startgettasklog() {
      if (this.logquery.name === "") {
        return;
      }
      gettaskLog(this.logquery).then(resp => {
        this.data = resp.data;
        this.pagecount = resp.count;
      });
    },
    handleNodeClick(data) {
      if (data.id !== undefined) {
        this.tasklog = data.taskresp_data;
        var tasktype = {
          1: "主任务",
          2: "父任务",
          3: "子任务"
        };

        this.tasklogtitle = `${tasktype[data.tasktype]} ${data.name}[${
          data.id
        }]`;
      } else {
        this.tasklog = "";
        this.tasklogtitle = "";
      }
    },
    renderContent(h, { node, data, store }) {
      return (
        <span class="custom-tree-node">
          <svg-icon icon-class={data.status} />
          <span style="margin-left:5px;">{data.name}</span>
        </span>
      );
    },
    gettasklogpre(task) {
      this.diatasktitle = `任务运行日志 ${task.name}`;
      this.diaogVisible = true;
      var treelog = {
        id: task.runby_taskid,
        start_time: task.start_time
      };
      gettaskLogTree(treelog).then(resp => {
        this.treetaskdata = resp.data;
      });
    },
    startcleantask() {
      var reqdata = {
        name: this.logquery.name,
        preday: this.cleanpredat
      };
      cleantasklog(reqdata).then(resp => {
        if (resp.code === 0) {
          Message.success(
            `总共清理任务 ${this.logquery.name} 的日志共 ${resp.data.delcount} 条`
          );
          this.cleanlogvisible = false;
          this.startgettasklog();
        } else {
          Message.error(`清理任务失败 ${resp.msg}`);
        }
      });
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
  color: #abadaf;
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
.el-table__empty-block {
  min-height: 60px;
  text-align: center;
  width: 100%;
  display: -webkit-box;
  display: -ms-flexbox;
  display: none;
  -webkit-box-pack: center;
  -ms-flex-pack: center;
  justify-content: center;
  -webkit-box-align: center;
  -ms-flex-align: center;
  align-items: center;
}
</style>


