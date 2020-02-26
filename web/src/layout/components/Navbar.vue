<template>
  <div class="navbar">
    <hamburger
      :is-active="sidebar.opened"
      class="hamburger-container"
      @toggleClick="toggleSideBar"
    />

    <breadcrumb class="breadcrumb-container" />

    <div class="right-menu">
      <el-badge :value="notifycount" class="badge">
        <router-link to="/notify">
          <i style="font-size: 16px;height:32px;" class="el-icon-bell"></i>
        </router-link>
      </el-badge>
      <el-dropdown size="small" placement="top">
        <div class="name">{{name}}</div>
        <el-dropdown-menu>
          <router-link to="/">
            <el-dropdown-item>首页</el-dropdown-item>
          </router-link>
          <router-link to="/profile">
            <el-dropdown-item>个人设置</el-dropdown-item>
          </router-link>
          <el-dropdown-item>
            <span style="display:block;" @click="logout">退出登录</span>
          </el-dropdown-item>
        </el-dropdown-menu>
      </el-dropdown>
    </div>
  </div>
</template>

<script>
import { mapGetters } from "vuex";
import Breadcrumb from "@/components/Breadcrumb";
import Hamburger from "@/components/Hamburger";
import { getnotify } from "@/api/notify";
export default {
  components: {
    Breadcrumb,
    Hamburger
  },
  computed: {
    ...mapGetters(["sidebar"])
  },
  created() {
    this.startgetnotifys();
    this.interval = setInterval(this.startgetnotifys, 5000);
  },
  data() {
    return {
      name: this.$store.getters.name,
      notifycount: 0,
      interval: null
    };
  },
  methods: {
    toggleSideBar() {
      this.$store.dispatch("app/toggleSideBar");
    },
    async logout() {
      await this.$store.dispatch("user/logout");
      window.clearInterval(this.interval);
      this.$router.push(`/login?redirect=${this.$route.fullPath}`);
    },
    startgetnotifys() {
      getnotify().then(resp => {
        this.notifycount = resp.data.length;
      });
    }
  }
};
</script>

<style lang="scss" scoped>
.navbar {
  height: 50px;
  overflow: hidden;
  position: relative;
  background: #ffffff;
  box-shadow: 0 1px 4px rgba(0, 21, 41, 0.08);

  .hamburger-container {
    line-height: 46px;
    height: 100%;
    float: left;
    cursor: pointer;
    transition: background 0.3s;
    -webkit-tap-highlight-color: transparent;

    &:hover {
      background: rgba(0, 0, 0, 0.025);
    }
  }

  .breadcrumb-container {
    float: left;
  }

  .right-menu {
    float: right;
    margin-right: 10px;
  }
  .badge {
    margin-right: 20px;
    margin-top: 13px;
  }
  .name {
    margin-right: 10px;
    margin-top: 15px;
    font-size: 15px;
    font-weight: bold;
  }
}
</style>
