<template>
  <div class="app-container">
    <div v-if="is_create === true || is_change == true" style="margin-left:25px;margin-right:80px">
      <el-form
        :model="task"
        ref="task"
        :rules="rules"
        label-position="right"
        label-width="120px"
        size="mini"
      >
        <el-form-item label="任务名称" prop="name">
          <el-input
            :disabled="is_change || is_preview"
            v-model="task.name"
            placeholder="请输入任务名称(建议任务名称格式为:标签_名字,方便后续搜索查找)"
            clearable
            style="width: 500px;"
            maxlength="30"
            show-word-limit
          ></el-input>
        </el-form-item>
        <el-form-item label="任务类型" prop="task_type" placeholder="请选择任务类型">
          <el-select
            :disabled="is_preview"
            clearable
            v-model="task.task_type"
            @change="changetaskdata"
          >
            <el-option
              v-for="item in task_typeoption"
              :key="item.label"
              :label="item.label"
              :value="item.value"
            ></el-option>
          </el-select>
        </el-form-item>
        <div style="margin-left:80px;margin-right:40px">
          <div v-if="task.task_type === 1">
            <ul>
              <el-select
                :disabled="is_preview"
                size="mini"
                v-model="savecode.lang"
                @change="changelang"
              >
                <el-option
                  v-for="item in langoption"
                  :key="item.label"
                  :label="item.label"
                  :value="item.value"
                ></el-option>
              </el-select>
              <div style="margin-top:5px;">
                <el-card :body-style="{ padding: '0px' }">
                  <editor
                    v-model="savecode.code"
                    theme="solarized_dark"
                    :lang="lang[savecode.lang]"
                    height="300"
                    @init="codeinitEditor"
                    :options="{ 
                      readOnly: is_preview,
                      wrap: 'free',
                      indentedSoftWrap: false,
                      enableBasicAutocompletion: true,
                      enableSnippets: true,
                      enableLiveAutocompletion: true
                 }"
                  ></editor>
                </el-card>
                <span style="text-align: right;color: #909399;font-size: 13px;">编辑框太小? 鼠标双击试试🤔</span>
                <br />
              </div>
            </ul>
          </div>
          <div v-if="task.task_type === 2">
            <ul>
              <div>
                <el-table
                  size="mini"
                  :data="methodurl"
                  style="width: 100%"
                  empty-text="please press add new"
                >
                  <el-table-column label="Method" min-width="15px;">
                    <template slot-scope="scope">
                      <el-select
                        :disabled="is_preview"
                        size="mini"
                        v-model="saveapi.method"
                        placeholder="请选择"
                      >
                        <el-option
                          v-for="item in methodoption"
                          :key="item.label"
                          :label="item.label"
                          :value="item.value"
                        ></el-option>
                      </el-select>
                    </template>
                  </el-table-column>
                  <el-table-column label="URL" min-width="100px;">
                    <template slot-scope="scope">
                      <el-input :disabled="is_preview" size="mini" v-model="saveapi.url"></el-input>
                    </template>
                  </el-table-column>
                </el-table>
              </div>
              <div>
                <el-table size="mini" :data="headerlist" style="width: 100%" empty-text=" ">
                  <el-table-column prop="name" label="Header List" min-width="100px;">
                    <template slot-scope="scope">
                      <el-input
                        :disabled="is_preview"
                        size="mini"
                        v-model="headerlist[scope.$index].key"
                      ></el-input>
                    </template>
                  </el-table-column>
                  <el-table-column prop="age" min-width="100px;">
                    <template slot-scope="scope">
                      <el-input
                        :disabled="is_preview"
                        size="mini"
                        v-model="headerlist[scope.$index].value"
                      ></el-input>
                    </template>
                  </el-table-column>
                  <el-table-column width="70px;">
                    <template slot-scope="scope">
                      <el-button
                        :disabled="is_preview"
                        size="mini"
                        @click="deleteheaderRow(scope.$index)"
                        icon="el-icon-delete"
                        type="info"
                        circle
                      ></el-button>
                    </template>
                  </el-table-column>
                </el-table>
                <div style="margin-left:11px">
                  <el-button
                    :disabled="is_preview"
                    type="info"
                    size="mini"
                    @click="addheader"
                  >Add New</el-button>
                </div>
              </div>
              <div v-if="['GET','HEAD'].includes(saveapi.method) === false">
                <div>
                  <el-table
                    size="mini"
                    :data="content_typetable"
                    style="width: 100%"
                    empty-text=" "
                  >
                    <el-table-column prop="name" label="Content Type" min-width="100px;">
                      <template slot-scope="scope">
                        <el-select
                          :disabled="is_preview"
                          size="mini"
                          filterable
                          v-model="content_type"
                          @change="edit_header"
                          style="width:100%;"
                        >
                          <el-option
                            v-for="item in content_typeoption"
                            :key="item.label"
                            :label="item.label"
                            :value="item.value"
                          ></el-option>
                        </el-select>
                      </template>
                    </el-table-column>
                  </el-table>
                </div>
                <div v-if="content_type == 'application/x-www-form-urlencoded'">
                  <el-table size="mini" :data="formlist" style="width: 100%" empty-text=" ">
                    <el-table-column prop="name" label="Form List" min-width="100px;">
                      <template slot-scope="scope">
                        <el-input
                          :disabled="is_preview"
                          size="mini"
                          v-model="formlist[scope.$index].key"
                        ></el-input>
                      </template>
                    </el-table-column>
                    <el-table-column prop="age" min-width="100px;">
                      <template slot-scope="scope">
                        <el-input
                          :disabled="is_preview"
                          size="mini"
                          v-model="formlist[scope.$index].value"
                        ></el-input>
                      </template>
                    </el-table-column>
                    <el-table-column width="70px;">
                      <template slot-scope="scope">
                        <el-button
                          :disabled="is_preview"
                          size="mini"
                          @click="deleteformRow(scope.$index)"
                          icon="el-icon-delete"
                          type="info"
                          circle
                        ></el-button>
                      </template>
                    </el-table-column>
                  </el-table>
                  <div style="margin-left:11px">
                    <el-button
                      :disabled="is_preview"
                      type="info"
                      size="mini"
                      @click="addform"
                    >Add New</el-button>
                  </div>
                </div>
                <div v-if="content_type == 'application/json'">
                  <el-table size="mini" :data="bodylist" style="width: 100%" empty-text=" ">
                    <el-table-column prop="name" label="Raw Requesy Body" min-width="100px;">
                      <!-- <div style="margin-top:5px;"> -->
                      <el-card :body-style="{ padding: '0px' }">
                        <editor
                          v-model="jsonplayload"
                          theme="solarized_dark"
                          lang="json"
                          height="300"
                          @init="rowreqbodyinitEditor"
                          :options="{ 
                            readOnly: is_preview,
                            wrap: 'free',
                            indentedSoftWrap: false}"
                        ></editor>
                      </el-card>
                    </el-table-column>
                  </el-table>
                </div>
              </div>
            </ul>
          </div>
        </div>
        <el-row>
          <el-col :span="11">
            <el-form-item label="父任务" prop="parent_taskids">
              <el-select
                :disabled="is_preview"
                multiple
                filterable
                v-model="task.parent_taskids"
                multiple-limit="20"
              >
                <el-option
                  v-for="item in taskselect"
                  :key="item.label"
                  :label="item.label"
                  :value="item.value"
                ></el-option>
              </el-select>
            </el-form-item>
          </el-col>
          <el-col :span="11">
            <el-form-item
              v-if="task.parent_taskids.length != 0"
              label="父任务运行策略"
              prop="parent_runparallel"
            >
              <el-select :disabled="is_preview" v-model="task.parent_runparallel">
                <el-option
                  v-for="item in parallelrun"
                  :key="item.label"
                  :label="item.label"
                  :value="item.value"
                ></el-option>
              </el-select>
            </el-form-item>
          </el-col>
        </el-row>
        <el-row>
          <el-col :span="11">
            <el-form-item label="子任务" prop="child_taskids">
              <el-select
                :disabled="is_preview"
                multiple
                filterable
                v-model="task.child_taskids"
                multiple-limit="20"
              >
                <el-option
                  v-for="item in taskselect"
                  :key="item.label"
                  :label="item.label"
                  :value="item.value"
                ></el-option>
              </el-select>
            </el-form-item>
          </el-col>
          <el-col :span="11">
            <el-form-item
              v-if="task.child_taskids.length != 0"
              label="子任务运行策略"
              prop="child_runparallel"
            >
              <el-select :disabled="is_preview" v-model="task.child_runparallel">
                <el-option
                  v-for="item in parallelrun"
                  :key="item.label"
                  :label="item.label"
                  :value="item.value"
                ></el-option>
              </el-select>
            </el-form-item>
          </el-col>
        </el-row>
        <el-form-item label="CronExpr" prop="cronexpr">
          <div>
            <el-input
              :disabled="is_preview"
              clearable
              v-model="task.cronexpr"
              placeholder="请输入Cron表达式"
              style="width: 500px;"
            >
              <template slot="append">
                <el-popover placement="top" width="650" v-model="cronPopover">
                  <cron v-show="cronPopover" v-model="cronExpression"></cron>
                  <div style="text-align: center; margin: 0">
                    <el-button type="primary" size="medium" @click="changeCron">确定</el-button>
                    <el-button type="warning" size="medium" @click="cronPopover = false">取消</el-button>
                  </div>
                  <el-button :disabled="is_preview" slot="reference">编辑</el-button>
                </el-popover>
              </template>
            </el-input>
          </div>
        </el-form-item>
        <el-form-item v-if="is_change" label="调度状态" prop="run">
          <el-switch
            :disabled="is_preview"
            v-model="task.run"
            active-text="正常调度"
            inactive-text="停止调度"
          ></el-switch>
        </el-form-item>
        <el-form-item label="超时时间(s)" prop="timeout">
          <el-input-number
            :disabled="is_preview"
            controls-position="right"
            v-model="task.timeout"
            :min="-1"
            label="超时时间"
          ></el-input-number>
        </el-form-item>
        <el-form-item label="路由策略" prop="route_policy">
          <el-select :disabled="is_preview" v-model="task.route_policy">
            <el-option
              v-for="item in route_policyoption"
              :key="item.label"
              :label="item.label"
              :value="item.value"
            ></el-option>
          </el-select>
        </el-form-item>
        <el-form-item label="报警策略" prop="alarm_status">
          <el-select :disabled="is_preview" v-model="task.alarm_status">
            <el-option
              v-for="item in alarm_statusoption"
              :key="item.value"
              :label="item.label"
              :value="item.value"
            ></el-option>
          </el-select>
        </el-form-item>
        <el-form-item label="报警用户" prop="alarm_userids">
          <el-select
            :disabled="is_preview"
            multiple
            filterable
            v-model="task.alarm_userids"
            multiple-limit="10"
          >
            <el-option
              v-for="item in userselect"
              :key="item.label"
              :label="item.label"
              :value="item.value"
            ></el-option>
          </el-select>
        </el-form-item>
        <el-form-item label="主机组" prop="host_groupid">
          <el-select :disabled="is_preview" filterable v-model="task.host_groupid">
            <el-option
              v-for="item in hostgroupselect"
              :key="item.label"
              :label="item.label"
              :value="item.value"
            ></el-option>
          </el-select>
        </el-form-item>
        <el-form-item label="Expect Code">
          <el-input
            :disabled="is_preview"
            type="number"
            v-model="task.expect_code"
            placeholder="期望返回码"
            clearable
            style="width: 500px;"
          ></el-input>
        </el-form-item>
        <el-form-item label="Expect Content">
          <el-input
            v-if="is_preview"
            :disabled="is_preview"
            type="textarea"
            v-model="task.expect_content"
            placeholder
            clearable
            style="width: 500px;"
          ></el-input>
          <el-input
            v-else
            type="textarea"
            v-model="task.expect_content"
            placeholder="期望返回内容"
            clearable
            style="width: 500px;"
          ></el-input>
        </el-form-item>
        <el-form-item label="备注">
          <el-input
            v-if="is_preview"
            :disabled="is_preview"
            type="textarea"
            v-model="task.remark"
            placeholder
            clearable
            style="width: 500px;"
            maxlength="100"
            show-word-limit
          ></el-input>
          <el-input
            v-else
            type="textarea"
            v-model="task.remark"
            placeholder="请输入备注"
            clearable
            style="width: 500px;"
            maxlength="100"
            show-word-limit
          ></el-input>
        </el-form-item>
        <el-form-item v-if="is_preview" label="创建时间">
          <el-input :disabled="is_preview" v-model="taskcreatetime" clearable></el-input>
        </el-form-item>
        <el-form-item v-if="is_preview" label="更新时间">
          <el-input :disabled="is_preview" v-model="taskupdatetime" clearable></el-input>
        </el-form-item>
      </el-form>
      <div style="margin-left: 120px;">
        <el-button
          v-if="is_preview === false"
          size="small"
          type="primary"
          @click="submittask('task')"
        >确 定</el-button>
        <el-button
          v-if="is_preview === false"
          size="small"
          @click="is_create = false;is_change = false"
        >取 消</el-button>
        <el-button v-if="is_preview" size="small" @click="is_preview=false;is_change=false">返 回</el-button>
      </div>
      <!-- </ul> -->
    </div>
    <div v-else-if="is_running === true">
      <!-- Running Task -->
      <div style="float: right">
        <el-tooltip class="item" effect="dark" content="任务列表" placement="top-start">
          <el-button type="warning" size="small" @click="toTask">Task</el-button>
        </el-tooltip>
      </div>
      <div>
        <!-- {{ runningdata }} -->
        <el-table
          v-loading="listLoading"
          :data="runningdata"
          stripe
          fit
          highlight-current-row
          style="width: 100%;"
        >
          <el-table-column align="center" label="任务名称">
            <template slot-scope="scope">
              <span>{{ scope.row.name }}</span>
            </template>
          </el-table-column>
          <el-table-column align="center" label="CronExpr">
            <template slot-scope="scope">
              <span>{{ scope.row.cronexpr }}</span>
            </template>
          </el-table-column>
          <el-table-column align="center" label="调度时间">
            <template slot-scope="scope">
              <span>{{ scope.row.start_timestr }}</span>
            </template>
          </el-table-column>
          <el-table-column align="center" label="运行时长">
            <template slot-scope="scope">
              <span>{{ scope.row.run_time / 1000 }}秒</span>
            </template>
          </el-table-column>
          <el-table-column align="center" label="触发方式">
            <template slot-scope="scope">
              <span>{{ scope.row.triggerstr }}</span>
            </template>
          </el-table-column>
          <el-table-column fixed="right" align="center" label="操作">
            <template slot-scope="scope">
              <el-button-group>
                <el-button type="warning" size="mini" @click="startopen(scope.row)">实时日志</el-button>
                <el-popconfirm
                  :hideIcon="true"
                  title="确定终止任务?"
                  @onConfirm="startkilltask(scope.row)"
                >
                  <el-button slot="reference" type="danger" size="mini">终止任务</el-button>
                </el-popconfirm>
              </el-button-group>
            </template>
          </el-table-column>
        </el-table>
        <div style="margin-top: 10px;float:right;height: 70px;">
          <el-pagination
            :page-size="runningquery.limit"
            @current-change="handleCurrentChangeruntask"
            background
            layout="total,prev, pager, next"
            :total="pagecount"
          ></el-pagination>
        </div>
        <el-dialog
          :title="diarealtasktitle"
          :visible.sync="diarealogVisible"
          center
          width="80%"
          @close="startclose"
        >
          <el-container>
            <el-aside width="200px">
              <el-tree
                highlight-current
                :default-expand-all="true"
                :data="runtaskdata"
                :render-content="renderContent"
                @node-click="handleNodeClick"
              ></el-tree>
            </el-aside>
            <!-- <el-divider direction="vertical" style="height: 150em;"></el-divider> -->
            <el-main>
              <span class="sub-title" v-html="realtasklogtitle"></span>
              <br />
              <span
                class="sub-title"
                v-if="realtasklogtitle !== '' && realtasklog === ''"
              >正在获取任务实时运行日志，请等待...</span>
              <el-card :body-style="{ padding: '0px' }">
                <editor
                  v-if="realtasklog != ''"
                  v-model="realtasklog"
                  theme="solarized_dark"
                  lang="text"
                  height="500"
                  width="100%"
                  @init="realloginitEditor"
                  :options="{ readOnly: true, wrap: 'free' ,vScrollBarAlwaysVisible: true,
    wrapBehavioursEnabled: true,
    autoScrollEditorIntoView: true,}"
                ></editor>
              </el-card>
            </el-main>
          </el-container>
        </el-dialog>
      </div>
    </div>
    <div v-else>
      <div style="float: right">
        <el-form :inline="true" label-width="80px">
          <el-form-item label="任务名称">
            <el-input
              v-model="query.psname"
              size="small"
              @keyup.enter.native="getalltask"
              placeholder="前置匹配模糊搜索"
              style="width:200px;margin-right: 1px"
            ></el-input>
          </el-form-item>
          <el-form-item>
            <el-tooltip
              v-if="query.self === false"
              class="item"
              effect="dark"
              content="查看自已创建的任务"
              placement="top-start"
            >
              <el-button
                v-if="query.self === false"
                type="success"
                size="small"
                @click="query.self = true;getalltask()"
              >Self</el-button>
            </el-tooltip>
            <el-tooltip
              v-if="query.self === true"
              class="item"
              effect="dark"
              content="查看所有任务"
              placement="top-start"
            >
              <el-button
                v-if="query.self === true"
                type="success"
                size="small"
                @click="query.self = false;getalltask()"
              >AllTask</el-button>
            </el-tooltip>
          </el-form-item>
          <el-form-item>
            <el-tooltip class="item" effect="dark" content="新建任务" placement="top-start">
              <el-button type="danger" size="small" @click="createtaskpre">New</el-button>
            </el-tooltip>
          </el-form-item>
          <el-form-item>
            <el-tooltip class="item" effect="dark" content="查看正在运行的任务" placement="top-start">
              <el-button type="warning" size="small" @click="getrunningtaskpre">Running</el-button>
            </el-tooltip>
          </el-form-item>
        </el-form>
      </div>
      <el-table
        v-loading="listLoading"
        :data="data"
        stripe
        fit
        highlight-current-row
        style="width: 100%;"
      >
        <el-table-column align="center" fixed="left" label="任务名称" min-width="100">
          <template slot-scope="scope">
            <span>{{ scope.row.name }}</span>
          </template>
        </el-table-column>
        <el-table-column align="center" label="任务类型" width="150">
          <template slot-scope="scope">
            <span v-if="scope.row.task_type == 1">
              <el-tag size="small" type="warning">{{ scope.row.task_typedesc }}</el-tag>
              <el-tag size="small" type="danger">{{ scope.row.task_data.langdesc }}</el-tag>
            </span>
            <span v-else>
              <el-tag size="small" type="warning">{{ scope.row.task_typedesc }}</el-tag>
              <el-tag size="small" type="danger">{{scope.row.task_data.method}}</el-tag>
            </span>
          </template>
        </el-table-column>
        <el-table-column align="center" label="路由策略" width="80">
          <template slot-scope="scope">
            <span>{{ scope.row.route_policydesc }}</span>
          </template>
        </el-table-column>
        <el-table-column align="center" label="CronExpr" min-width="100">
          <template slot-scope="scope">
            <span>{{ scope.row.cronexpr }}</span>
          </template>
        </el-table-column>

        <el-table-column align="center" label="调度状态" width="80">
          <template slot-scope="scope">
            <el-tag v-if="scope.row.run" size="mini" type="success">Normal</el-tag>
            <el-tag v-else size="mini" type="danger">Stop</el-tag>
          </template>
        </el-table-column>

        <el-table-column align="center" label="超时(s)" width="70">
          <template slot-scope="scope">
            <span v-if="scope.row.timeout === -1">-</span>
            <span v-else>{{ scope.row.timeout }}</span>
          </template>
        </el-table-column>

        <el-table-column align="center" label="主机组" min-width="100">
          <template slot-scope="scope">
            <span>{{ scope.row.host_group }}</span>
          </template>
        </el-table-column>
        <el-table-column align="center" label="报警策略" min-width="70">
          <template slot-scope="scope">
            <span>{{ scope.row.alarm_statusdesc }}</span>
          </template>
        </el-table-column>

        <el-table-column align="center" label="创建人" width="80">
          <template slot-scope="scope">
            <span>{{ scope.row.create_by }}</span>
          </template>
        </el-table-column>

        <el-table-column align="center" label="备注" min-width="70">
          <template slot-scope="scope">
            <span>{{ scope.row.remark }}</span>
          </template>
        </el-table-column>
        <el-table-column fixed="right" align="center" label="操作" width="130">
          <template slot-scope="scope">
            <el-button-group>
              <el-button type="danger" size="mini" @click="changetaskpre(scope.row)">修改</el-button>
              <el-tooltip class="item" effect="dark" content="从此任务克隆一个新的任务" placement="top-start">
                <el-button
                  type="warning"
                  size="mini"
                  @click="clonevisible=true;clonename=scope.row.name;cloneid=scope.row.id;clonenewname='';"
                >克隆</el-button>
              </el-tooltip>
              <el-popconfirm
                :hideIcon="true"
                title="删除任务后对应的任务日志也会被删除，确定删除任务?"
                @onConfirm="startdeletetask(scope.row)"
              >
                <el-button slot="reference" type="danger" size="mini">删除</el-button>
              </el-popconfirm>
            </el-button-group>
            <br />
            <el-button-group>
              <el-button type="info" size="mini" @click="previewtaskpre(scope.row)">详情</el-button>
              <el-button type="success" size="mini">
                <router-link :to="{name: 'Log', query: {name: scope.row.name}}">日志</router-link>
              </el-button>
              <el-popconfirm :hideIcon="true" title="立即运行任务?" @onConfirm="startruntask(scope.row)">
                <el-button type="primary" slot="reference" size="mini">运行</el-button>
              </el-popconfirm>
            </el-button-group>
          </template>
        </el-table-column>
      </el-table>
      <div style="margin-top: 10px;float:right;height: 70px;">
        <el-pagination
          :page-size="query.limit"
          @current-change="handleCurrentChangetask"
          background
          layout="total,prev, pager, next"
          :total="pagecount"
        ></el-pagination>
      </div>
      <el-dialog :title="clonediatitle" center :visible.sync="clonevisible" width="20%">
        <el-input
          v-model="clonenewname"
          size="mini"
          placeholder="请输入新的任务名称"
          style="width: 100%;"
          maxlength="30"
          show-word-limit
        ></el-input>
        <p></p>
        <div style="text-align: center;">
          <el-button type="primary" size="mini" @click="startclonetask">确定</el-button>
          <el-button size="mini" type="text" @click="clonevisible = false">取消</el-button>
        </div>
      </el-dialog>
    </div>
  </div>
</template>

<script>
import {
  gettask,
  createtask,
  changetask,
  deletetask,
  killtask,
  runtask,
  getrunningtasks,
  gettaskLog,
  getselecttask,
  clonetask
} from "@/api/task";

import { getToken } from "@/utils/auth";
import { gethostgroup, getselecthostgroup } from "@/api/hostgroup";

import { getselectuser } from "@/api/user";

import store from "@/store";
import { isNumber } from "@/utils/validate";
import { Message } from "element-ui";

import cron from "@/components/Cron/cron";

export default {
  components: {
    editor: require("vue2-ace-editor"),
    cron
  },
  computed: {
    CronCollapse() {
      return `点击查看【${this.cronExpression}】最近五次的运行时间`;
    },
    clonediatitle() {
      return `克隆任务 ${this.clonename}`;
    }
  },
  watch: {
    "task.cronexpr": {
      handler(newv, oldv) {
        this.cronExpression = newv;
      },
      deep: true,
      immediate: true
    },
    diarealogVisible: {
      handler(newv, oldv) {
        if (newv === false && this.tlsocket != null) {
          console.log("start close socket");
          this.tlsocket.close();
        }
      }
    }
  },
  data() {
    return {
      listLoading: false,
      cronexprdialog: false,
      createjobdialog: false,
      task: {
        id: "",
        name: "",
        task_type: "",
        task_data: {},
        parent_taskids: [],
        parent_runparallel: false,
        child_taskids: [],
        child_runparallel: false,
        host_groupid: "",
        timeout: -1,
        run: false,
        cronexpr: "",
        alarm_userids: [],
        route_policy: 1,
        expect_code: 0,
        expect_content: "",
        alarm_status: -1,
        remark: ""
      },
      status_options: [
        {
          value: false,
          label: "正常"
        },
        {
          value: true,
          label: "停止"
        }
      ],
      task_typeoption: [
        {
          value: 1,
          label: "Code"
        },
        {
          value: 2,
          label: "API"
        }
      ],
      parallelrun: [
        {
          value: true,
          label: "并行"
        },
        {
          value: false,
          label: "串行"
        }
      ],
      route_policyoption: [
        {
          value: 1,
          label: "Random"
        },
        {
          value: 2,
          label: "RoundRobin"
        },
        {
          value: 3,
          label: "Weight"
        },
        {
          value: 4,
          label: "LeastTask"
        }
      ],
      alarm_statusoption: [
        {
          value: -2,
          label: "Always"
        },
        {
          value: -1,
          label: "Fail"
        },
        {
          value: 1,
          label: "Success"
        }
      ],
      data: [],
      runningdata: [],
      pagecount: 0,
      rules: {
        name: [{ required: true, message: "请输入任务名称", trigger: "blur" }],
        task_type: [
          { required: true, message: "请选择任务类型", trigger: "blur" }
        ],
        cronexpr: [
          { required: true, message: "请输入Cron表达式", trigger: "blur" }
        ],
        timeout: [
          { required: true, message: "请输入超时时间", trigger: "change" }
        ],
        route_policy: [
          { required: true, message: "请选择路由策略", trigger: "blur" }
        ],
        alarm_status: [
          { required: true, message: "请选择报警策略", trigger: "blur" }
        ],
        alarm_userids: [
          { required: true, message: "请选择报警用户", trigger: "blur" }
        ],
        host_groupid: [
          { required: true, message: "请选择主机组", trigger: "blur" }
        ]
      },
      is_change: false,
      is_create: false,
      is_running: false,
      is_preview: false,
      savecode: {
        lang: 1,
        code: ""
      },
      langoption: [
        {
          value: 1,
          label: "shell"
        },
        {
          value: 4,
          label: "python"
        },
        {
          value: 2,
          label: "python3"
        },
        {
          value: 3,
          label: "golang"
        },
        {
          value: 5,
          label: "nodejs"
        }
      ],
      lang: {
        1: "sh",
        2: "python3",
        3: "golang",
        4: "python",
        5: "nodejs"
      },
      saveapi: {
        url: "",
        method: "GET",
        payload: "",
        header: {}
      },
      method: "",
      url: "",
      methodoption: [
        {
          value: "GET",
          label: "GET"
        },
        {
          value: "HEAD",
          label: "HEAD"
        },
        {
          value: "POST",
          label: "POST"
        },
        {
          value: "PUT",
          label: "PUT"
        },
        {
          value: "PATCH",
          label: "PATCH"
        },
        {
          value: "DELETE",
          label: "DELETE"
        }
      ],
      content_type: "",
      content_typeoption: [
        {
          value: "application/json",
          label: "application/json"
        },
        {
          value: "application/x-www-form-urlencoded",
          label: "application/x-www-form-urlencoded"
        }
      ],
      headerlist: [{}],
      methodurl: [{}],
      content_typetable: [{}],
      formlist: [{}],
      bodylist: [{}],
      cronPopover: false,
      cronExpression: "",
      jsonplayload: "",
      query: {
        offset: 0,
        limit: 15,
        self: false,
        psname: ""
      },
      runningquery: {
        offset: 0,
        limit: 15
      },
      examplecode: {
        1: `#!/usr/bin/env sh
function main() {
    echo "run shell"
}

main
        `,
        2: `#!/usr/bin/env python3
def main():
    print("run python3")

if __name__ == '__main__':
    main()`,
        4: `#!/usr/bin/env python
def main():
    print("run python")

if __name__ == '__main__':
    main()`,
        3: `package main

import "fmt"

func main() {
	fmt.Println("run golang")
}`,
        5: `#!/usr/bin/env node
console.log("run nodejs")`
      },
      hostgroupselect: [],
      userselect: [],
      taskselect: [],
      dialogVisible: false,
      diarealogVisible: false,
      trsocket: null, // task run status
      tlsocket: null, // task run log,
      runtaskdata: [],
      diarealtasktitle: "",
      diatasktitle: "",
      runningInterval: undefined,
      realtasklog: "",
      realtasklogtitle: "",
      currenttasklogid: "",
      taskcreatetime: "",
      taskupdatetime: "",
      selftask: true,
      clonenewname: "",
      cloneid: "",
      clonename: "",
      clonevisible: false
    };
  },
  created() {
    this.getalltask();
  },
  methods: {
    getalltask() {
      gettask(this.query).then(resp => {
        this.data = resp.data;
        this.pagecount = resp.count;
      });
    },
    submittask(formName) {
      this.$refs[formName].validate(valid => {
        if (valid) {
          this.task.expect_code = parseInt(this.task.expect_code);
          if (isNaN(this.task.expect_code)) {
            Message.error("Execpt Code only allow is number");
          }
          if (this.task.task_type === 1) {
            this.task.task_data = this.savecode;
          } else if (this.task.task_type === 2) {
            // add header
            this.headerlist.forEach(item => {
              if (item.key !== "") {
                this.saveapi.header[item.key] = item.value;
              }
            });
            // add method
            if (["GET", "HEAD"].includes(this.saveapi.method) === false) {
              if (this.content_type === "application/x-www-form-urlencoded") {
                var formlist = [];
                this.formlist.forEach(item => {
                  if (item.key !== "") {
                    formlist.push(`${item.key}=${item.value}`);
                  }
                });
                this.saveapi.payload = formlist.join("&");
              } else if (this.content_type === "application/json") {
                this.saveapi.payload = this.jsonplayload;
              }
            }
            this.task.task_data = this.saveapi;
          } else {
            console.log("err: support task type", this.task.task_type);
          }

          if (this.is_create === true) {
            delete this.task.id;
            createtask(this.task).then(response => {
              if (response.code === 0) {
                Message.success(`创建任务 ${this.task.name} 成功`);
                this.getalltask();
                this.is_create = false;
              } else {
                Message.error(
                  `创建任务 ${this.task.name} 失败: ${response.msg}`
                );
              }
            });
          } else if (this.is_change === true) {
            var name = this.task.name;

            delete this.task.task_data.langdesc;
            changetask(this.task).then(response => {
              if (response.code === 0) {
                Message.success(`修改任务 ${name} 成功`);
                this.getalltask();
                this.is_change = false;
              } else {
                Message.error(`修改任务 ${name} 失败: ${response.msg}`);
              }
            });
          }
        } else {
          return false;
        }
      });
    },
    startdeletetask(task) {
      var deldata = {
        id: task.id
      };
      deletetask(deldata).then(response => {
        if (response.code === 0) {
          this.getalltask();
          Message.success(`删除任务${task.name}成功`);
        } else {
          Message.error(`删除任务${task.name}失败: ${response.msg}`);
        }
      });
    },
    startkilltask(task) {
      var killdata = {
        id: task.id
      };
      killtask(killdata).then(response => {
        if (response.code === 0) {
          Message.success(`终止任务 ${task.name} 成功`);
        } else {
          Message.error(`终止任务 ${task.name} 失败: ${response.msg}`);
        }
      });
    },
    startruntask(task) {
      var rundata = {
        id: task.id
      };
      runtask(rundata).then(response => {
        if (response.code === 0) {
          Message.success(`任务 ${task.name} 已经开始运行`);
        } else {
          Message.error(`运行任务${task.name}失败: ${response.msg}`);
        }
      });
    },
    changelang() {
      this.savecode.code = this.examplecode[this.savecode.lang];
    },
    edit_header() {
      this.saveapi.header["Content-Type"] = this.content_type;
    },
    codeinitEditor: function(editor) {
      editor.setAutoScrollEditorIntoView(true);
      editor.setShowPrintMargin(false);
      editor.on("dblclick", function() {
        editor.container.webkitRequestFullscreen();
      });
      require("brace/ext/language_tools");
      require("brace/mode/sh");
      require("brace/mode/python");
      require("brace/mode/javascript");
      // require("brace/mode/text");
      // require("brace/mode/json");
      require("brace/mode/golang");
      require("brace/theme/solarized_dark");
    },
    realloginitEditor: function(editor) {
      editor.setAutoScrollEditorIntoView(true);
      editor.setShowPrintMargin(false);
      editor.on("change", function() {
        editor.renderer.scrollToLine(Number.POSITIVE_INFINITY);
      });
      // require("brace/ext/language_tools");
      // require("brace/mode/sh");
      // require("brace/mode/python");
      // require("brace/mode/javascript");
      require("brace/mode/text");
      // require("brace/mode/json");
      // require("brace/mode/golang");
      require("brace/theme/solarized_dark");
    },
    rowreqbodyinitEditor: function(editor) {
      editor.setAutoScrollEditorIntoView(true);
      editor.setShowPrintMargin(false);

      // require("brace/ext/language_tools");
      // require("brace/mode/sh");
      // require("brace/mode/python");
      // require("brace/mode/javascript");
      // require("brace/mode/text");
      require("brace/mode/json");
      // require("brace/mode/golang");
      require("brace/theme/solarized_dark");
    },
    addheader() {
      this.headerlist.push({});
    },
    deleteheaderRow(index) {
      this.headerlist.splice(index, 1);
    },
    addform() {
      this.formlist.push({});
    },
    deleteformRow(index) {
      this.formlist.splice(index, 1);
    },
    generatedict(datalist) {},
    changeCron() {
      this.task.cronexpr = this.cronExpression;
      this.cronPopover = false;
    },
    changetaskdata() {
      if (this.task.task_type === 1) {
        if (this.is_create === true) {
          this.task.expect_code = 0;
          this.savecode.code = this.examplecode[this.savecode.lang];
        }
        this.task.task_data = this.savecode;
      } else if (this.task.task_type === 2) {
        if (this.is_create === true) {
          this.task.expect_code = 200;
        }
        this.task.task_data = this.saveapi;
      }
    },
    createtaskpre() {
      if (this.task.name == undefined) {
        this.tasl["name"] = "";
      } else {
        this.task.name = "";
      }
      this.task.run = true;
      this.task.task_type = "";
      this.task.task_data = {};
      this.task.parent_taskids = [];
      // this.task.parent_runparallel = task.parent_runparallel;
      this.task.child_taskids = [];
      // this.task.child_runparallel = task.child_runparallel;
      this.task.host_groupid = [];
      this.task.timeout = -1;
      this.task.cronexpr = "";
      this.task.alarm_userids = [];
      this.task.route_policy = 1;
      this.task.expect_code = 0;
      this.task.expect_content = "";
      this.task.alarm_status = -1;
      this.task.remark = "";

      this.is_create = true;
      this.savecode = { lang: 1, code: "" };
      this.saveapi = { url: "", method: "GET", payload: "", header: {} };
      this.jsonplayload = "";

      this.headerlist = [{}];
      this.methodurl = [{}];
      this.content_typetabl = [{}];
      this.formlist = [{}];
      this.bodylist = [{}];

      this.getselect();
    },
    getselect() {
      getselecttask().then(response => {
        this.taskselect = response.data;
      });
      getselecthostgroup().then(response => {
        this.hostgroupselect = response.data;
      });
      getselectuser().then(response => {
        this.userselect = response.data;
      });
    },
    changetaskpre(task) {
      this.getselect();

      if (this.task.id === undefined) {
        this.task["id"] = task.id;
      } else {
        this.task.id = task.id;
      }
      this.task.name = task.name;
      this.task.run = task.run;
      this.task.task_type = task.task_type;
      this.task_data = task.task_data;
      this.task.parent_taskids = task.parent_taskids;
      this.task.parent_runparallel = task.parent_runparallel;
      this.task.child_taskids = task.child_taskids;
      this.task.child_runparallel = task.child_runparallel;
      this.task.host_groupid = task.host_groupid;
      this.task.timeout = task.timeout;
      this.task.cronexpr = task.cronexpr;
      this.task.alarm_userids = task.alarm_userids;
      this.task.route_policy = task.route_policy;
      this.task.expect_code = task.expect_code;
      this.task.expect_content = task.expect_content;
      this.task.alarm_status = task.alarm_status;
      this.task.remark = task.remark;

      this.is_change = true;

      if (this.task.task_type === 1) {
        this.savecode = task.task_data;
      }
      if (this.task.task_type === 2) {
        this.saveapi = task.task_data;
        this.headerlist = [];
        for (let key in this.saveapi.header) {
          if (key !== "Content-Type") {
            this.headerlist.push({
              key: key,
              value: this.saveapi.header[key]
            });
          }
        }
        if (this.headerlist.length === 0) {
          this.headerlist.push({});
        }
        this.content_type = this.saveapi.header["Content-Type"];

        if (this.content_type == "application/json") {
          this.jsonplayload = this.saveapi.payload;
        }
        if (this.content_type === "application/x-www-form-urlencoded") {
          // testetet=tetetstet&tetste=12121
          this.formlist = [];
          if (this.saveapi.payload !== "") {
            this.saveapi.payload.split("&").forEach(item => {
              var v = item.split("=");
              this.formlist.push({
                key: v[0],
                value: v[1]
              });
            });
          }
          if (this.formlist.length === 0) {
            this.formlist.push({});
          }
        }
      }
    },
    previewtaskpre(task) {
      this.is_preview = true;
      this.taskcreatetime = task.create_time;
      this.taskupdatetime = task.update_time;
      this.changetaskpre(task);
    },
    getrunningtaskpre(query) {
      this.is_running = true;

      this.getrundatafunc(query);

      this.runningInterval = setInterval(this.getrundatafunc, 5000);
    },
    getrundatafunc(query) {
      getrunningtasks(query).then(resp => {
        this.runningdata = resp.data;
        this.pagecount = resp.count;
      });
    },
    handleCurrentChangetask(val) {
      this.query.offset = (val - 1) * this.query.limit;
      this.getalltask(this.query);
    },
    handleCurrentChangeruntask(val) {
      this.runningquery.offset = (val - 1) * this.runningquery.limit;
      this.getrunningtaskpre();
    },
    startopen(task) {
      this.realtasklogtitle = "";
      this.realtasklog = "";
      this.currenttasklogid = task.id;
      this.diarealtasktitle = `实时任务日志: ${task.name}`;
      this.diarealogVisible = true;
      this.runtaskdata = [];
      var token = getToken();
      var host = "";
      if (process.env.NODE_ENV === "production") {
        host = `${location.origin}`;
      } else {
        const config = require("../../../vue.config.js");
        host = config.devServer.proxy[process.env.VUE_APP_BASE_API].target;
      }
      var wsurl = `${host.replace(
        "http",
        "ws"
      )}/api/v1/task/status/websocket?id=${task.id}`;

      this.trsocket = new WebSocket(wsurl);

      this.trsocket.onopen = event => {
        this.trsocket.send(token);
      };
      this.trsocket.onmessage = event => {
        this.runtaskdata = JSON.parse(event.data);
        this.trsocket.send("ok");
      };
    },
    startclose() {
      if (this.trsocket !== null) {
        // this.trsocket.send("close");
        console.log(`start close websocket task run status`);
        this.trsocket.close();
      }
      if (this.tlsocket !== null) {
        // this.tlsocket.send("close");
        console.log(`start close websocket task log status`);
        this.tlsocket.close();
      }
    },
    handleNodeClick(data, node, label) {
      if (data.id == undefined) {
        return;
      }
      if (this.tlsocket !== null) {
        console.log("start close websocket");
        this.tlsocket.close();
      }
      var tasktype = {
        1: "主任务",
        2: "父任务",
        3: "子任务"
      };
      this.realtasklogtitle = `${tasktype[data.tasktype]} ${data.name}[${
        data.id
      }]`;

      this.realtasklog = "";
      var token = getToken();

      var host = "";
      if (process.env.NODE_ENV === "production") {
        host = `${location.origin}`;
      } else {
        const config = require("../../../vue.config.js");
        host = config.devServer.proxy[process.env.VUE_APP_BASE_API].target;
      }
      var wsurl = `${host.replace("http", "ws")}/api/v1/task/log/websocket?id=${
        this.currenttasklogid
      }&realid=${data.id}&type=${data.tasktype}`;
      console.log(`start conn websocket ${wsurl}`);
      this.tlsocket = new WebSocket(wsurl);
      console.log(this.tlsocket);
      this.tlsocket.onopen = event => {
        this.tlsocket.send(token);
      };
      this.tlsocket.onmessage = event => {
        // this.runtaskdata = JSON.parse(event.data);
        this.realtasklog = this.realtasklog + event.data;
        this.tlsocket.send("ok");
      };
    },
    renderContent(h, { node, data, store }) {
      return (
        <span class="custom-tree-node">
          <svg-icon icon-class={data.status} />
          <span style="margin-left:5px;">{data.name}</span>
        </span>
      );
    },
    toTask() {
      this.is_running = false;
      window.clearInterval(this.runningInterval);
    },

    startclonetask() {
      var reqdata = {
        id: this.cloneid,
        name: this.clonenewname
      };
      if (this.clonenewname === "") {
        Message.warning("请输入新的任务名称");
        return;
      }
      clonetask(reqdata).then(resp => {
        if (resp.code === 0) {
          Message.success(`克隆任务成功`);
          this.clonevisible = false;
          this.getalltask();
        } else {
          Message.error(`克隆任务失败 ${resp.msg}`);
        }
      });
    }
  }
};
</script>

<style lang="scss" scoped>
.el-button--mini {
  padding: 5px 5px;
}
.sub-title {
  text-align: left;
  color: #909399;
  font-size: 16px;
  font-family: "Helvetica Neue, Helvetica, PingFang SC, Hiragino Sans GB, Microsoft YaHei, Arial, sans-serif";
  // margin-bottom: 6px;
  font-weight: 700;
}
</style>
