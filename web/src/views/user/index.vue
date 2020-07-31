<template>
  <div class="app-container">
    <div v-if="is_change || is_create" style="margin-left:25px;margin-right:80px；height:80px">
      <el-form
        :model="user"
        ref="user"
        :rules="rules"
        label-position="right"
        label-width="120px"
        size="mini"
      >
        <el-form-item label="用户名" prop="name">
          <el-input
            :disabled="is_change"
            v-model="user.name"
            placeholder="请输入用户名"
            clearable
            style="width: 500px;"
            maxlength="30"
            show-word-limit
          ></el-input>
        </el-form-item>
        <el-form-item label="用户类型" prop="role">
          <el-select v-model="user.role">
            <el-option
              v-for="item in roleoptions"
              :key="item.label"
              :label="item.label"
              :value="item.value"
            ></el-option>
          </el-select>
        </el-form-item>
        <el-form-item label="密码" prop="password">
          <el-tooltip content="修改密码" placement="top">
            <el-switch v-if="!is_create" v-model="changepasswd"></el-switch>
          </el-tooltip>
          <span v-if="changepasswd ">
            <br />
          </span>
          <el-input
            v-if="changepasswd || is_create"
            v-model="password1"
            type="password"
            placeholder="请输入密码"
            clearable
            style="width: 500px;"
          ></el-input>
          <span v-if="changepasswd || is_create">
            <br />
          </span>
          <el-input
            v-if="changepasswd || is_create"
            type="password"
            v-model="password2"
            placeholder="请再次输入密码"
            clearable
            style="width: 500px;"
          ></el-input>
        </el-form-item>
        <el-form-item v-if="is_change" label="状态" prop="forbid">
          <el-select v-model="user.forbid">
            <el-option
              v-for="item in forbidoptions"
              :key="item.label"
              :label="item.label"
              :value="item.value"
            ></el-option>
          </el-select>
        </el-form-item>
        <el-form-item label="备注" prop="remark">
          <el-input
            type="textarea"
            v-model="user.remark"
            placeholder="请输入任务备注"
            clearable
            style="width: 500px;"
          ></el-input>
        </el-form-item>
      </el-form>
      <div style="margin-left: 120px;">
        <el-button size="small" type="primary" @click="submituser('user')">确 定</el-button>
        <el-button size="small" @click="is_create = false;is_change = false">取 消</el-button>
      </div>
    </div>
    <div v-else>
      <div style="float: right">
        <el-tooltip class="item" effect="dark" content="创建用户" placement="top-start">
          <el-button type="primary" size="small" @click="createuserpre">New</el-button>
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
        <el-table-column align="center" fixed="left" label="用户名" min-width="100">
          <template slot-scope="scope">
            <span>{{ scope.row.name }}</span>
          </template>
        </el-table-column>

        <el-table-column align="center" label="用户类型" min-width="70">
          <template slot-scope="scope">
            <el-tag v-if="scope.row.role === 1" type="warning" size="mini">{{ scope.row.rolestr }}</el-tag>
            <el-tag v-if="scope.row.role === 2" type="danger" size="mini">{{ scope.row.rolestr }}</el-tag>
            <el-tag v-if="scope.row.role === 3" type="info" size="mini">{{ scope.row.rolestr }}</el-tag>
          </template>
        </el-table-column>

        <el-table-column align="center" label="状态" min-width="70">
          <template slot-scope="scope">
            <el-switch :value="!scope.row.forbid" active-color="#13ce66" inactive-color="#ff4949"></el-switch>
          </template>
        </el-table-column>
        <el-table-column property="create_time" label="创建时间" width="160"></el-table-column>
        <el-table-column align="center" label="备注" min-width="50">
          <template slot-scope="scope">
            <span>{{ scope.row.remark }}</span>
          </template>
        </el-table-column>
        <el-table-column fixed="right" align="center" label="操作" min-width="50">
          <template slot-scope="scope">
            <el-button-group>
              <el-button type="warning" size="mini" @click="changeuserpre(scope.row)">修改</el-button>
              <el-popconfirm
                :hideIcon="true"
                title="确定删除此用户(只能删除非管理员用户)"
                @onConfirm="deleteuser(scope.row.id)"
              >
                <el-button slot="reference" type="danger" size="mini">删除</el-button>
              </el-popconfirm>
            </el-button-group>
          </template>
        </el-table-column>
      </el-table>
      <div style="margin-top: 10px;float:right;height: 70px;">
        <el-pagination
          :page-size="userquery.limit"
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
  getallusers,
  adminchangeinfo,
  createuser,
  admindeleteuser,
} from "@/api/user";

import { Message } from "element-ui";

export default {
  data() {
    return {
      listLoading: false,
      data: [],
      pagecount: 0,
      rules: {
        name: [{ required: true, message: "请输入用户名称", trigger: "blur" }],
        role: [{ required: true, message: "请选择用户类型", trigger: "blur" }],
        password: [
          { required: true, message: "请输入用户密码", trigger: "blur" },
        ],
        forbid: [
          { required: true, message: "请输入用户密码", trigger: "blur" },
        ],
      },
      addips: [],
      pagecount: 0,
      userquery: {
        offset: 0,
        limit: 15,
      },
      user: {
        id: "",
        name: "",
        role: 0,
        forbid: 1,
        password: "",
        remark: "",
      },
      createuser: {
        name: "",
        password: "",
        role: 0,
        remark: "",
      },
      roleoptions: [
        {
          label: "普通用户",
          value: 1,
        },
        {
          label: "管理员",
          value: 2,
        },
        {
          label: "访客",
          value: 3,
        },
      ],
      forbidoptions: [
        {
          label: "禁止登陆",
          value: true,
        },
        {
          label: "正常登陆",
          value: false,
        },
      ],
      hghosts: [],
      is_change: false,
      is_create: false,
      hostselect: [],
      changepasswd: false,
      password1: "",
      password2: "",
    };
  },
  created() {
    this.startgetallusers();
  },
  methods: {
    startgetallusers() {
      getallusers(this.userquery).then((response) => {
        this.data = response.data;
        this.pagecount = response.count;
      });
    },
    createuserpre() {
      if (this.user.name == undefined) {
        this.user["name"] = "";
      }
      this.user.name = "";
      this.user.password = "";
      this.user.role = 1;
      this.user.remark = "";
      this.is_create = true;
      this.password1 === "";
      this.password2 === "";
    },
    changeuserpre(user) {
      if (this.user.id == undefined) {
        this.user["id"] = "";
      }
      if (this.user.forbid == undefined) {
        this.user["forbid"] = 1;
      }
      if (user.forbid === 1) {
        this.forbid = true;
      } else if (user.forbid === 2) {
        this.forbid = false;
      }
      this.changepasswd = false;
      this.user.id = user.id;
      this.user.name = user.name;
      this.user.role = user.role;
      this.user.forbid = user.forbid;

      this.user.remark = user.remark;
      this.user.password = "";
      this.is_change = true;
      this.password1 === "";
      this.password2 === "";
    },

    submituser(formName) {
      if (this.changepasswd || this.is_create) {
        if (this.password1 === "" || this.password2 === "") {
          Message.warning("密码不能为空，请重新输入");
          return;
        }
        if (this.password1 !== this.password2) {
          Message.warning("两次输入密码不一致，请重新输入");
          return;
        }
        try {
          window.btoa(`${this.password1}`);
        } catch (error) {
          Message.warning("密码只能使用字母、数字、符号");
          return;
        }
        if (this.password1.length < 8) {
          Message.warning("密码最少8位");
          return;
        }
        this.user.password = this.password1;
      }
      try {
        window.btoa(`${this.user.name}`);
      } catch (error) {
        Message.warning("用户名只能使用字母、数字、符号");
        return;
      }
      if (this.is_change === true) {
        var name = this.user.name;
        delete this.user.name;
        if (this.changepasswd === false) {
          this.user.password = "";
        }
        adminchangeinfo(this.user).then((resp) => {
          if (resp.code === 0) {
            Message.success(`修改用户 ${name} 成功`);
            this.startgetallusers();
            this.is_change = false;
          } else {
            Message.error(`修改用户 ${name} 失败: ${resp.msg}`);
          }
        });
      } else {
        if (this.password1 === "" || this.password2 === "") {
          Message.warning("密码不能为空，请重新输入");
          return;
        }
        if (this.password1 !== this.password2) {
          Message.warning("两次输入密码不一致，请重新输入");
          return;
        }
        try {
          window.btoa(`${this.password1}`);
        } catch (error) {
          Message.warning("密码只能使用字母、数字、符号");
          return;
        }

        this.$refs[formName].validate((valid) => {
          if (valid) {
            if (this.is_create === true) {
              // if (this.password1 === "" || this.password2 === "") {
              //   Message.error("密码不能为空，请重新输入");
              //   return;
              // }
              // if (this.password1 !== this.password2) {
              //   Message.error("两次输入密码不一致，请重新输入");
              //   return;
              // }
              // try {
              //   window.btoa(`${this.user.name}:${this.password1}`);
              // } catch (error) {
              //   Message.error("用户名和密码只能使用字母、数字、符号");
              //   return;
              // }

              // this.user.password = this.password1;
              delete this.user.id;
              delete this.user.forbid;

              createuser(this.user).then((resp) => {
                if (resp.code === 0) {
                  Message.success(`创建用户 ${this.user.name} 成功`);
                  this.startgetallusers();
                  this.is_create = false;
                } else {
                  Message.error(`创建用户 ${name} 失败: ${resp.msg}`);
                }
              });
            }
          }
        });
      }
    },
    deleteuser(id) {
      admindeleteuser({ id: id }).then((resp) => {
        if (resp.code === 0) {
          Message.success(`删除用户 ${this.user.name} 成功`);
          this.startgetallusers();
          this.is_create = false;
        } else {
          Message.error(`删除用户 ${name} 失败: ${resp.msg}`);
        }
      });
    },
    handleCurrentChangerun(page) {
      this.userquery.offset = (page - 1) * this.userquery.limit;
      this.startgetallusers();
    },
  },
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
