// cron 组件来自 https://gitee.com/lindeyi/vue-cron

<template>
  <div class="cron" :val="value_">
    <el-row>
      <el-col :span="18">
        <el-tabs v-model="activeName">
          <el-tab-pane label="秒" name="s">
            <second-and-minute v-model="sVal" lable="秒"></second-and-minute>
          </el-tab-pane>
          <el-tab-pane label="分" name="m">
            <second-and-minute v-model="mVal" lable="分"></second-and-minute>
          </el-tab-pane>
          <el-tab-pane label="时" name="h">
            <hour v-model="hVal" lable="时"></hour>
          </el-tab-pane>
          <el-tab-pane label="日" name="d">
            <day v-model="dVal" lable="日" :disable="disableday"></day>
          </el-tab-pane>
          <el-tab-pane label="月" name="month">
            <month v-model="monthVal" lable="月"></month>
          </el-tab-pane>
          <el-tab-pane label="周" name="week">
            <week v-model="weekVal" lable="周" :disable="disableweek"></week>
          </el-tab-pane>
          <el-tab-pane label="年" name="year">
            <year v-model="yearVal" lable="年"></year>
          </el-tab-pane>
        </el-tabs>
        <el-table :data="tableData" size="mini" border style="width: 100%;">
          <el-table-column prop="sVal" label="秒" min-width="60"></el-table-column>
          <el-table-column prop="mVal" label="分" min-width="60"></el-table-column>
          <el-table-column prop="hVal" label="时" min-width="60"></el-table-column>
          <el-table-column prop="dVal" label="日" min-width="60"></el-table-column>
          <el-table-column prop="monthVal" label="月" min-width="60"></el-table-column>
          <el-table-column prop="weekVal" label="周" min-width="60"></el-table-column>
          <el-table-column prop="yearVal" label="年" min-width="60"></el-table-column>
        </el-table>
      </el-col>
      <el-col :span="6">
        <el-collapse v-model="collapseactiveName" @change="handleChange">
          <el-collapse-item title="最近运行时间" name="1">
            <div v-for="item in cronrecenttime" :key="item">
              <span>{{ item }}</span>
            </div>
          </el-collapse-item>
        </el-collapse>
      </el-col>
    </el-row>
  </div>
</template>

<script>
import SecondAndMinute from "./secondAndMinute";
import hour from "./hour";
import day from "./day";
import month from "./month";
import week from "./week";
import year from "./year";
import { parsecron } from "@/api/task";
export default {
  props: {
    value: {
      type: String
    }
  },
  data() {
    return {
      activeName: "s",
      collapseactiveName: ["1"],
      sVal: "",
      mVal: "",
      hVal: "",
      dVal: "",
      monthVal: "",
      weekVal: "",
      yearVal: "",
      cronrecenttime: [],
      disableweek: false,
      disableday: false
    };
  },
  watch: {
    value(a, b) {
      this.updateVal();
      this.handleChange();
    }
  },
  computed: {
    tableData() {
      return [
        {
          sVal: this.sVal,
          mVal: this.mVal,
          hVal: this.hVal,
          dVal: this.dVal,
          monthVal: this.monthVal,
          weekVal: this.weekVal,
          yearVal: this.yearVal
        }
      ];
    },
    value_() {
      if (!this.dVal && !this.weekVal) {
        return "";
      }
      // if (this.dVal === "?" && this.weekVal === "?") {
      //   this.$message.error("日期与星期不可以同时为“不指定”");
      // }
      // if (this.dVal !== "?" && this.weekVal !== "?") {
      //   this.$message.error("日期与星期必须有一个为“不指定”");
      // }
      if (this.dVal === "?") {
        this.disableweek = false;
      } else {
        this.disableweek = true;
      }

      if (this.weekVal === "?") {
        this.disableday = false;
      } else {
        this.disableday = true;
      }
      // 日 不指定  周 可以选
      // 日 指定 周 禁止
      // 周 不指定 日 可以选
      // 周 指定 日 禁止选

      let v = `${this.sVal} ${this.mVal} ${this.hVal} ${this.dVal} ${this.monthVal} ${this.weekVal} ${this.yearVal}`;
      if (v !== this.value) {
        this.$emit("input", v);
      }
      return v;
    }
  },
  methods: {
    updateVal() {
      if (!this.value) {
        return;
      }
      let arrays = this.value.split(" ");
      this.sVal = arrays[0];
      this.mVal = arrays[1];
      this.hVal = arrays[2];
      this.dVal = arrays[3];
      this.monthVal = arrays[4];
      this.weekVal = arrays[5];
      this.yearVal = arrays[6];
    },
    handleChange() {
      if (this.value === "") {
        return;
      }
      var query = {
        expr: window.btoa(this.value)
      };
      parsecron(query).then(response => {
        this.cronrecenttime = response.data;
      });
    }
  },
  created() {
    this.updateVal();
    this.handleChange()
  },
  components: {
    SecondAndMinute,
    hour,
    day,
    month,
    week,
    year
  }
};
</script>

<style lang="css">
.cron {
  text-align: left;
  padding: 10px;
  background: #fff;
  border: 1px solid #dcdfe6;
  box-shadow: 0 2px 4px 0 rgba(0, 0, 0, 0.12), 0 0 6px 0 rgba(0, 0, 0, 0.04);
}
</style>
