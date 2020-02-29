<template>
  <div class="app-container">
    <div style="margin-left:25px;margin-right:80px">
      <!-- <el-tabs v-model="activename" @tab-click="handleClick">
        <el-tab-pane name="setting">
          <span slot="label">
            <svg-icon icon-class="usersetting" />设置
          </span> -->
          <el-form label-position="right" label-width="80px" :model="userinfo">
            <el-form-item label="修改密码">
              <el-tooltip content="修改密码" placement="top">
                <el-switch v-model="changepasswd"></el-switch>
              </el-tooltip>
              <span v-if="changepasswd">
                <br />
              </span>
              <el-input
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
              ></el-input>
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
            <!-- <el-form-item label="Slack">
              <el-input
                :disabled="!alarmstatus.slack"
                v-model="userinfo.slack"
                size="mini"
                style="width: 300px;"
              ></el-input>
            </el-form-item> -->
            <el-form-item label="Telegram">
              <el-input
                :disabled="!alarmstatus.telegram"
                v-model="userinfo.telegram"
                size="mini"
                style="width: 300px;"
              ></el-input>
            </el-form-item>
            <el-form-item label="备注">
              <el-input type="textarea" v-model="userinfo.remark" size="mini" style="width: 300px;"></el-input>
            </el-form-item>
          </el-form>
          <div style="margin-left: 80px;">
            <el-popconfirm
              :hideIcon="true"
              title="确定修改个人信息?"
              @onConfirm="submitchangeinfo"
            >
              <el-button slot="reference" size="small" type="primary">更 新</el-button>
            </el-popconfirm>
          </div>
        <!-- </el-tab-pane> -->
        <!-- <el-tab-pane name="operatelog">
          <span slot="label">
            <svg-icon icon-class="operate" />操作日志
          </span>
          <div class="block">
            <el-timeline>
              <el-timeline-item timestamp="2018/4/12" placement="top">
                <el-card>
                  <h4>更新 Github 模板</h4>
                  <p>王小虎 提交于 2018/4/12 20:46</p>
                </el-card>
              </el-timeline-item>
              <el-timeline-item timestamp="2018/4/3" placement="top">
                <el-card>
                  <h4>更新 Github 模板</h4>
                  <p>王小虎 提交于 2018/4/3 20:46</p>
                </el-card>
              </el-timeline-item>
              <el-timeline-item timestamp="2018/4/2" placement="top">
                <el-card>
                  <h4>更新 Github 模板</h4>
                  <p>王小虎 提交于 2018/4/2 20:46</p>
                </el-card>
              </el-timeline-item>
            </el-timeline>
          </div>
        </el-tab-pane> -->
      </el-tabs>
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
      userinfo: {
        id: "",
        email: "",
        wechat: "",
        dingphone: "",
        slack: "",
        telegram: "",
        password: "",
        remark: ""
      },
      changepasswd: false,
      alarmstatus: {
        email: false,
        dingphone: false,
        slack: false,
        telegram: false,
        wechat: false,
        wehook: false
      }
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
      getInfo().then(resp => {
        var data = resp.data;
        this.userinfo.id = data.id;
        this.userinfo.email = data.email;
        this.userinfo.wechat = data.wechat;
        this.userinfo.dingphone = data.dingphone;
        this.userinfo.slack = data.slack;
        this.userinfo.telegram = data.telegram;
        this.userinfo.remark = data.remark;
      });
    },
    submitchangeinfo() {
      if (this.changepasswd) {
        if (
          this.password1 === "" ||
          this.password2 === "" ||
          this.password1 !== this.password2
        ) {
          Message.error("两次密码输入不一致请重新输入");
          return;
        } else {
          try {
            window.btoa(`${this.password1}`);
          } catch (error) {
            Message.error("密码只能使用字母、数字、符号");
            return;
          }
          this.userinfo.password = this.password1;
        }
      }
      changeselfinfo(this.userinfo).then(resp => {
        if (resp.code === 0) {
          Message.success("更新成功");
          this.changepasswd = false;
          this.password1 = "";
          this.password2 = "";
        } else {
          Message.error(`更新失败 ${resp.msg}`);
        }
      });
    },
    startgetalarmstatus() {
      getalarmstatus().then(resp => {
        this.alarmstatus = resp.data;
      });
    }
  }
};
</script>