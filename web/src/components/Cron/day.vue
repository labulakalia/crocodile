<template >
  <div :val="value_">
    <div>
      <el-radio :disabled="disable" v-model="type" label="1" size="mini" border>每日</el-radio>
    </div>
    <div>
      <el-radio :disabled="disable" v-model="type" label="5" size="mini" border>不指定</el-radio>
    </div>
    <div>
      <el-radio :disabled="disable" v-model="type" label="2" size="mini" border>周期</el-radio>
      <span style="margin-left: 0px; margin-right: 0px;">从</span>
      <el-input-number
        :disabled="disable"
        @change="type = '2'"
        v-model="cycle.start"
        :min="1"
        :max="31"
        size="mini"
        style="width: 100px;"
      ></el-input-number>
      <span style="margin-left: 0px; margin-right: 0px;">至</span>
      <el-input-number
        :disabled="disable"
        @change="type = '2'"
        v-model="cycle.end"
        :min="2"
        :max="31"
        size="mini"
        style="width: 100px;"
      ></el-input-number>日
    </div>
    <div>
      <el-radio :disabled="disable" v-model="type" label="3" size="mini" border>循环</el-radio>
      <span style="margin-left: 0px; margin-right: 0px;">从</span>
      <el-input-number
        :disabled="disable"
        @change="type = '3'"
        v-model="loop.start"
        :min="1"
        :max="31"
        size="mini"
        style="width: 100px;"
      ></el-input-number>
      <span style="margin-left: 0px; margin-right: 0px;">日开始，每</span>
      <el-input-number
        :disabled="disable"
        @change="type = '3'"
        v-model="loop.end"
        :min="1"
        :max="31"
        size="mini"
        style="width: 100px;"
      ></el-input-number>日执行一次
    </div>
    <div>
      <el-radio :disabled="disable" v-model="type" label="8" size="mini" border>工作日</el-radio>
      <span style="margin-left: 0px; margin-right: 0px;">本月</span>
      <el-input-number
        :disabled="disable"
        @change="type = '8'"
        v-model="work"
        :min="1"
        :max="7"
        size="mini"
        style="width: 100px;"
      ></el-input-number>号，最近的工作日
    </div>
    <div>
      <el-radio :disabled="disable" v-model="type" label="6" size="mini" border>本月最后一天</el-radio>
    </div>
    <div>
      <el-radio :disabled="disable" v-model="type" label="4" size="mini" border>指定</el-radio>
      <el-checkbox-group v-model="appoint" :min="1">
        <div v-for="i in 6" :key="i">
          <el-checkbox
            :disabled="disable"
            @change="type = '4'"
            v-for="j in 6"
            v-if="(i-1) * 6 + j < 32"
            :key="j"
            :label="check((i-1) * 6 + j)"
          ></el-checkbox>
        </div>
      </el-checkbox-group>
    </div>
  </div>
</template>

<script>
export default {
  props: {
    value: {
      type: String,
      default: "?"
    },
    disable: {
      type: Boolean
    }
  },
  data() {
    return {
      type: "5", // 类型
      cycle: {
        // 周期
        start: 0,
        end: 0
      },
      loop: {
        // 循环
        start: 0,
        end: 0
      },
      week: {
        // 指定周
        start: 0,
        end: 0
      },
      work: 0,
      last: 0,
      appoint: [] // 指定
    };
  },
  computed: {
    value_() {
      let result = [];
      switch (this.type) {
        case "1": // 每秒
          result.push("*");
          break;
        case "2": // 周期
          result.push(`${this.cycle.start}-${this.cycle.end}`);
          break;
        case "3": // 循环
          result.push(`${this.loop.start}/${this.loop.end}`);
          break;
        case "4": // 指定
          result.push(this.appoint.join(","));
          break;
        case "6": // 最后
          result.push(`${this.last === 0 ? "" : this.last}L`);
          break;
        case "7": // 指定周
          result.push(`${this.week.start}#${this.week.end}`);
          break;
        case "8": // 工作日
          result.push(`${this.work}W`);
          break;
        default:
          // 不指定
          result.push("?");
          break;
      }
      this.$emit("input", result.join(""));
      return result.join("");
    }
  },
  watch: {
    value(a, b) {
      this.updateVal();
    }
  },
  methods: {
    updateVal() {
      if (!this.value) {
        return;
      }
      if (this.value === "?") {
        this.type = "5";
      } else if (this.value.indexOf("-") !== -1) {
        // 2周期
        if (this.value.split("-").length === 2) {
          this.type = "2";
          this.cycle.start = this.value.split("-")[0];
          this.cycle.end = this.value.split("-")[1];
        }
      } else if (this.value.indexOf("/") !== -1) {
        // 3循环
        if (this.value.split("/").length === 2) {
          this.type = "3";
          this.loop.start = this.value.split("/")[0];
          this.loop.end = this.value.split("/")[1];
        }
      } else if (this.value.indexOf("*") !== -1) {
        // 1每
        this.type = "1";
      } else if (this.value.indexOf("L") !== -1) {
        // 6最后
        this.type = "6";
        this.last = this.value.replace("L", "");
      } else if (this.value.indexOf("#") !== -1) {
        // 7指定周
        if (this.value.split("#").length === 2) {
          this.type = "7";
          this.week.start = this.value.split("#")[0];
          this.week.end = this.value.split("#")[1];
        }
      } else if (this.value.indexOf("W") !== -1) {
        // 8工作日
        this.type = "8";
        this.work = this.value.replace("W", "");
      } else {
        // *
        this.type = "4";
        this.appoint = this.value.split(",");
      }
    },
    check(val) {
      if (val < 10) {
        return "0" + val;
      }
      return val.toString();
    }
  },
  created() {
    this.updateVal();
  }
};
</script>

<style lang="css">
.el-checkbox + .el-checkbox {
  margin-left: 10px;
}
</style>
