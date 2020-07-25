<template>
  <div class="app-container">
    <div style="margin-left:25px;margin-right:80px">
      <!-- <el-tabs v-model="activename" @tab-click="handleClick">
        <el-tab-pane name="setting">
          <span slot="label">
            <svg-icon icon-class="usersetting" />设置
      </span>-->
      <el-form label-position="right" label-width="80px" :model="userinfo">
        <el-form-item label="修改密码">
          <el-tooltip content="修改密码" placement="top">
            <el-switch v-model="changepasswd"></el-switch>
          </el-tooltip>

          <el-form v-if="changepasswd" :model="pass" ref="pass" :rules="rules" size="mini">
            <el-form-item prop="password1">
              <el-input
                type="password"
                v-model="pass.password1"
                size="mini"
                clearable
                placeholder="请输入密码"
                style="width: 300px;"
              ></el-input>
            </el-form-item>
            <el-form-item prop="password2">
              <el-input
                type="password"
                v-model="pass.password2"
                size="mini"
                clearable
                placeholder="请再次输入密码"
                style="width: 300px;"
              ></el-input>
            </el-form-item>
          </el-form>

          <!-- <el-input
            v-if="changepasswd"
            type="password"
            v-model="password1"
            size="mini"
            clearable
            placeholder="请输入密码"
            style="width: 300px;"
          ></el-input>
          <span v-if="changepasswd">
            <br />
          </span>
          <el-input
            type="password"
            v-if="changepasswd"
            v-model="password2"
            size="mini"
            clearable
            placeholder="请再次输入密码"
            style="width: 300px;"
          ></el-input>-->
        </el-form-item>
        <el-form-item label="用户名">
          <el-input v-model="userinfo.name" size="mini" style="width: 300px;"></el-input>
        </el-form-item>
        <el-form-item label="邮箱">
          <el-input
            :disabled="!alarmstatus.email"
            v-model="userinfo.email"
            size="mini"
            style="width: 300px;"
          ></el-input>
        </el-form-item>
        <el-form-item label="WeChat">
          <el-input
            :disabled="!alarmstatus.wechat"
            v-model="userinfo.wechat"
            size="mini"
            style="width: 300px;"
          ></el-input>
        </el-form-item>
        <el-form-item label="钉钉">
          <el-input
            :disabled="!alarmstatus.dingphone"
            v-model="userinfo.dingphone"
            size="mini"
            style="width: 300px;"
          ></el-input>
        </el-form-item>
        <el-form-item label="Telegram">
          <el-input
            :disabled="!alarmstatus.telegram"
            v-model="userinfo.telegram"
            size="mini"
            style="width: 300px;"
          ></el-input>
        </el-form-item>
        <el-form-item label="备注">
          <el-input
            type="textarea"
            v-model="userinfo.remark"
            size="mini"
            style="width: 300px;"
            maxlength="100"
            show-word-limit
          ></el-input>
        </el-form-item>
      </el-form>
      <div style="margin-left: 80px;">
        <el-popconfirm :hideIcon="true" title="确定修改个人信息?" @onConfirm="submitchangeinfo">
          <el-button slot="reference" size="small" type="primary">更 新</el-button>
        </el-popconfirm>
      </div>
    </div>
  </div>
</template>

<script>
import { getInfo, changeselfinfo, getalarmstatus } from "@/api/user";
import { Message } from "element-ui";
export default {
  data() {
    return {
      activename: "setting",
      name: this.$store.getters.name,
      password1: "",
      password2: "",
      pass: {
        password1: "",
        password2: "",
      },
      rules: {
        password1: [{ required: true, message: "请输入密码", trigger: "blur" }],
        password2: [
          { required: true, message: "请再次输入密码", trigger: "blur" },
        ],
      },
      userinfo: {
        id: "",
        name: "",
        email: "",
        wechat: "",
        dingphone: "",
        slack: "",
        telegram: "",
        password: "",
        remark: "",
      },
      changepasswd: false,
      alarmstatus: {
        email: false,
        dingphone: false,
        slack: false,
        telegram: false,
        wechat: false,
        wehook: false,
      },
    };
  },
  created() {
    this.getuserinfo();
    this.startgetalarmstatus();
  },
  methods: {
    handleClick(tab, event) {
      console.log(tab);
    },
    getuserinfo() {
      getInfo().then((resp) => {
        var data = resp.data;
        this.userinfo.id = data.id;
        this.userinfo.email = data.email;
        this.userinfo.wechat = data.wechat;
        this.userinfo.dingphone = data.dingphone;
        this.userinfo.slack = data.slack;
        this.userinfo.telegram = data.telegram;
        this.userinfo.remark = data.remark;
        this.userinfo.name = data.name;
      });
    },
    submitchangeinfo() {
      if (this.changepasswd) {
        this.$refs["pass"].validate((valid) => {
          if (this.changepasswd && valid) {
            if (this.pass.password1 !== this.pass.password2) {
              Message.warning("两次密码输入不一致请重新输入");
              return;
            } else {
              try {
                window.btoa(`${this.pass.password1}`);
              } catch (error) {
                Message.warning("密码只能使用字母、数字、符号");
                return;
              }
              if (this.pass.password1.length < 8) {
                Message.warning("密码最少8位");
                return;
              }
              this.userinfo.password = this.pass.password1;
            }
          } else {
            return false;
          }
          changeselfinfo(this.userinfo).then((resp) => {
            if (resp.code === 0) {
              Message.success("更新成功");
              this.changepasswd = false;
              this.password1 = "";
              this.password2 = "";
            } else {
              Message.error(`更新失败 ${resp.msg}`);
            }
          });
        });
      } else {
        changeselfinfo(this.userinfo).then((resp) => {
          if (resp.code === 0) {
            Message.success("更新成功");
            this.changepasswd = false;
            this.password1 = "";
            this.password2 = "";
          } else {
            Message.error(`更新失败 ${resp.msg}`);
          }
        });
      }
      // this.$refs["pass"].validate(valid => {
      //   if (this.changepasswd && valid) {
      //     if (this.pass.password1 !== this.pass.password2) {
      //       Message.warning("两次密码输入不一致请重新输入");
      //       return;
      //     } else {
      //       try {
      //         window.btoa(`${this.pass.password1}`);
      //       } catch (error) {
      //         Message.warning("密码只能使用字母、数字、符号");
      //         return;
      //       }
      //       this.userinfo.password = this.pass.password1;
      //     }
      //   } else {
      //     return false;
      //   }
      //   changeselfinfo(this.userinfo).then(resp => {
      //     if (resp.code === 0) {
      //       Message.success("更新成功");
      //       this.changepasswd = false;
      //       this.password1 = "";
      //       this.password2 = "";
      //     } else {
      //       Message.error(`更新失败 ${resp.msg}`);
      //     }
      //   });
      // });

      // if (this.changepasswd) {
      //   this.$refs["pass"].validate(valid => {
      //     if (valid) {
      //       if (this.password1 !== this.password2) {
      //         Message.warning("两次密码输入不一致请重新输入");
      //         return;
      //       } else {
      //         try {
      //           window.btoa(`${this.password1}`);
      //         } catch (error) {
      //           Message.warning("密码只能使用字母、数字、符号");
      //           return;
      //         }
      //         this.userinfo.password = this.password1;
      //       }
      //       return true;
      //     } else {
      //       Message.warning("密码只能使用字母、数字、符号");
      //       return false;
      //     }
      //   });
      //   if (this.password1 === "" || this.password2 === "") {
      //     Message.warning("请输入密码");
      //     return;
      //   }
      // }
    },
    startgetalarmstatus() {
      getalarmstatus().then((resp) => {
        this.alarmstatus = resp.data;
      });
    },
  },
};
</script>