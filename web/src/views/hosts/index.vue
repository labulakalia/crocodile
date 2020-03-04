<template>
  <div class="app-container">
    <div style="margin-left:25px;margin-right:20px">
      <el-table :data="data">
        <el-table-column align="center" property="addr" label="IP" min-width="100"></el-table-column>
        <el-table-column align="center" property="hostname" label="主机名" min-width="100"></el-table-column>
        <el-table-column align="center" property="weight" label="权重" min-width="60"></el-table-column>
        <el-table-column align="center" property="version" label="版本" min-width="60"></el-table-column>
        <el-table-column align="center" label="状态" min-width="70">
          <template slot-scope="scope">
            <el-tag type="success" size="mini" v-if="scope.row.online">Online</el-tag>
            <el-tag type="danger" size="mini" v-else>Offline</el-tag>
          </template>
        </el-table-column>
        <el-table-column align="center" label="暂停" min-width="80">
          <template slot-scope="scope">
            <el-tag v-if="scope.row.stop" size="mini" type="danger">Stop</el-tag>
            <el-tag v-else size="mini" type="success">Normal</el-tag>
          </template>
        </el-table-column>
        <el-table-column align="center" label="操作" min-width="60">
          <template slot-scope="scope">
            <el-button-group>
              <el-tooltip
                v-if="scope.row.stop"
                class="item"
                effect="dark"
                content="恢复在此Worker上运行任务"
                placement="top"
              >
                <el-button type="success" size="mini" @click="changestate(scope.row.id)">正常</el-button>
              </el-tooltip>
              <el-tooltip v-else class="item" effect="dark" content="暂停此Worker运行任务" placement="top">
                <el-button type="warning" size="mini" @click="changestate(scope.row.id)">暂停</el-button>
              </el-tooltip>
              <el-popconfirm
                :hideIcon="true"
                title="删除此主机?"
                @onConfirm="startdeletehost(scope.row.id)"
              >
                <el-button type="danger" slot="reference" size="mini">删除</el-button>
              </el-popconfirm>
            </el-button-group>
          </template>
        </el-table-column>
      </el-table>
      <div style="margin-top: 10px;float:right;height: 70px;">
        <el-pagination
         :page-size="hostquery.limit"
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
import { gethost, stophost, deletehost } from "@/api/host";
import { Message } from "element-ui";
export default {
  data() {
    return {
      data: [],
      pagecount: 0,
      hostquery: {
        offset: 0,
        limit: 15
      }
    };
  },
  created() {
    this.startgethost();
  },
  methods: {
    startgethost() {
      gethost(this.hostquery).then(resp => {
        this.data = resp.data;
        this.pagecount = resp.count;
      });
    },
    changestate(id) {
      var data = {
        id: id,
      };
      stophost(data).then(resp => {
        if (resp.code === 0) {
          Message.success("修改状态成功");
          this.startgethost();
        } else {
          Message.error(`修改状态失败 errmsg: ${resp.msg}`);
        }
      });
    },
    handleCurrentChangerun(page) {
      this.hostquery.offset = (page - 1) * this.hostquery.limit;
      this.startgethost();
    },
    startdeletehost(id) {
      var deldata = {
        id: id
      };
      deletehost(deldata).then(resp => {
        if (resp.code === 0) {
          Message.success("删除成功");
          this.startgethost();
        } else {
          Message.warning(`删除失败 ${resp.msg}`);
        }
      });
    }
  }
};
</script>

<style>
.el-button--mini {
  padding: 5px 5px;
}
</style>