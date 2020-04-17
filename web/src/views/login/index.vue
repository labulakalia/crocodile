<template>
  <div
    class="login-container"
    v-loading="installloading"
    element-loading-text="正在安装..."
    element-loading-spinner="el-icon-loading"
    element-loading-background="rgba(0, 0, 0, 0.8)"
  >
    <el-form
      ref="loginForm"
      :model="loginForm"
      :rules="loginRules"
      class="login-form"
      auto-complete="on"
      label-position="left"
    >
      <div class="title-container">
        <h3 class="title">Crocodile分布式任务调度平台</h3>
        <h6 v-show="needinstall" class="installtitle">首次运行请先创建默认管理员用户然后进行安装操作</h6>
      </div>

      <el-form-item prop="username">
        <span class="svg-container">
          <svg-icon icon-class="user" />
        </span>
        <el-input
          ref="username"
          v-model="loginForm.username"
          placeholder="Username"
          name="username"
          type="text"
          tabindex="1"
          auto-complete="on"
          maxlength="30"
        />
      </el-form-item>

      <el-form-item prop="password">
        <span class="svg-container">
          <svg-icon icon-class="password" />
        </span>
        <el-input
          :key="passwordType"
          ref="password"
          v-model="loginForm.password"
          :type="passwordType"
          placeholder="Password"
          name="password"
          tabindex="2"
          auto-complete="on"
        />
        <!-- <span class="show-pwd" @click="showPwd(passwordType)">
          <svg-icon :icon-class="passwordType === 'password' ? 'eye' : 'eye-open'" />
        </span>-->
      </el-form-item>
      <el-form-item v-show="needinstall">
        <span class="svg-container">
          <svg-icon icon-class="password" />
        </span>
        <el-input
          :key="passwordType2"
          ref="password"
          v-model="password2"
          :type="passwordType2"
          placeholder="Password"
          name="password"
          tabindex="2"
          auto-complete="on"
        />
        <!-- <span class="show-pwd" @click="showPwd(passwordType2)">
          <svg-icon :icon-class="passwordType2 === 'password' ? 'eye' : 'eye-open'" />
        </span>-->
      </el-form-item>
      <el-button
        v-if="needinstall"
        :loading="loading"
        type="primary"
        style="width:100%;margin-bottom:30px;"
        @click.native.prevent="startinstallcrocodile"
      >开始安装</el-button>
      <el-button
        v-else
        :loading="loading"
        type="primary"
        style="width:100%;margin-bottom:30px;"
        @click.native.prevent="handleLogin"
      >登陆</el-button>
      <br />
    </el-form>
  </div>
</template>

<script>
import { validUsername } from "@/utils/validate";
import { queryinstallstatus, startinstall } from "@/api/install";
import { Message } from "element-ui";
import { login, logout } from "@/api/user";

export default {
  name: "Login",
  data() {
    const validateUsername = (rule, value, callback) => {
      if (!validUsername(value)) {
        callback(new Error("Please enter the correct user name"));
      } else {
        callback();
      }
    };
    const validatePassword = (rule, value, callback) => {
      if (value.length < 6) {
        callback(new Error("The password can not be less than 6 digits"));
      } else {
        callback();
      }
    };
    return {
      loginForm: {
        username: "",
        password: ""
      },
      password2: "",
      loginRules: {
        username: [
          { required: true, trigger: "blur", message: "请输入用户名" }
        ],
        password: [{ required: true, trigger: "blur", message: "请输入密码" }]
      },
      loading: false,
      passwordType: "password",
      passwordType2: "password",
      redirect: undefined,
      needinstall: false,
      installloading: false
    };
  },
  watch: {
    $route: {
      handler: function(route) {
        this.redirect = route.query && route.query.redirect;
      },
      immediate: true
    }
  },
  created() {
    this.startqueryinstallstatus();
  },
  methods: {
    startqueryinstallstatus() {
      queryinstallstatus().then(resp => {
        if (resp.code === 10700) {
          this.needinstall = true;
          Message.warning(resp.msg);
        }
      });
    },
    startinstallcrocodile() {
      // loginForm
      this.$refs["loginForm"].validate(valid => {
        if (valid) {
          if (this.loginForm.password !== this.password2) {
            Message.warning("两次密码输入不相同");
            return;
          }
          try {
            window.btoa(
              `${this.loginForm.username}:${this.loginForm.password}`
            );
          } catch (error) {
            Message.warning("用户名和密码只能使用字母、数字、符号");
            return;
          }
          if (this.loginForm.password.length < 8) {
            Message.warning("密码最少8位");
            return;
          }
          this.installloading = true;
          startinstall(this.loginForm)
            .then(resp => {
              if (resp.code === 0) {
                this.startqueryinstallstatus();
                this.installloading = false;
                this.needinstall = false;
                Message.success("恭喜你已经安装成功🎉");
              } else {
                Message.error(resp.msg);
                this.installloading = false;
                this.needinstall = false;
              }
            })
            .catch(err => {
              Message.error(err);
              console.log(err);
              this.installloading = false;
            });
        } else {
          return false;
        }
      });
    },
    handleLogin() {
      this.$refs.loginForm.validate(valid => {
        if (valid) {
          this.$store
            .dispatch("user/login", this.loginForm)
            .then(() => {
              this.$router.push({ path: this.redirect || "/" });
              this.loading = false;
            })
            .catch(() => {
              this.loading = false;
            });
        } else {
          return false;
        }
      });
    }
  }
};
</script>

<style lang="scss">
/* 修复input 背景不协调 和光标变色 */
/* Detail see https://github.com/PanJiaChen/vue-element-admin/pull/927 */

$bg: #283443;
$light_gray: #fff;
$cursor: #fff;

@supports (-webkit-mask: none) and (not (cater-color: $cursor)) {
  .login-container .el-input input {
    color: $cursor;
  }
}

/* reset element-ui css */
.login-container {
  .el-input {
    display: inline-block;
    height: 47px;
    width: 85%;

    input {
      background: transparent;
      border: 0px;
      -webkit-appearance: none;
      border-radius: 0px;
      padding: 12px 5px 12px 15px;
      color: $light_gray;
      height: 47px;
      caret-color: $cursor;

      &:-webkit-autofill {
        box-shadow: 0 0 0px 1000px $bg inset !important;
        -webkit-text-fill-color: $cursor !important;
      }
    }
  }

  .el-form-item {
    border: 1px solid rgba(255, 255, 255, 0.1);
    background: rgba(0, 0, 0, 0.1);
    border-radius: 5px;
    color: #454545;
  }
}
</style>

<style lang="scss" scoped>
$bg: #2d3a4b;
$dark_gray: #889aa4;
$light_gray: #eee;

.login-container {
  min-height: 100%;
  width: 100%;
  background-color: $bg;
  overflow: hidden;

  .login-form {
    position: relative;
    width: 520px;
    max-width: 100%;
    padding: 160px 35px 0;
    margin: 0 auto;
    overflow: hidden;
  }

  .tips {
    font-size: 14px;
    color: #fff;
    margin-bottom: 10px;

    span {
      &:first-of-type {
        margin-right: 16px;
      }
    }
  }

  .svg-container {
    padding: 6px 5px 6px 15px;
    color: $dark_gray;
    vertical-align: middle;
    width: 30px;
    display: inline-block;
  }

  .title-container {
    position: relative;

    .title {
      font-size: 26px;
      color: #bebcbc;
      margin: 0px auto 40px auto;
      text-align: center;
      font-weight: bold;
    }
    .installtitle {
      font-size: 16px;
      color: #f34747;
      margin: 0px auto 40px auto;
      text-align: center;
      font-weight: bold;
    }
  }

  .show-pwd {
    position: absolute;
    right: 10px;
    top: 7px;
    font-size: 16px;
    color: $dark_gray;
    cursor: pointer;
    user-select: none;
  }
}
</style>
