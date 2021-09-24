<template>
  <div>
    <v-tabs
      v-model="currentMode"
      grow
    >
      <v-tab>{{ $t('servers.Hourly') }}</v-tab>
      <v-tab>{{ $t('servers.Daily') }}</v-tab>
      <v-tab>{{ $t('servers.Weekly') }}</v-tab>
      <v-tab>{{ $t('servers.Monthly') }}</v-tab>
      <v-tab>{{ $t('servers.Advanced') }}</v-tab>
    </v-tabs>
    <v-tabs-items
      v-model="currentMode"
      class="pt-2"
    >
      <v-tab-item>
        <span class="tab">
          <i18n path="servers.HourlyTab">
            <template v-slot:hourInterval>
              <ui-input
                v-model="hour"
                dense
                class="cronDigit"
                @input="updateHour($event)"
              />
            </template>
            <template v-slot:minute>
              <ui-input
                v-model="minute"
                dense
                class="cronDigit"
                @input="updateMinute($event)"
              />
            </template>
          </i18n>
        </span>
      </v-tab-item>
      <v-tab-item>
        <span class="tab">
          <i18n path="servers.DailyTab">
            <template v-slot:dayInterval>
              <ui-input
                v-model="dayOfMonth"
                dense
                class="cronDigit"
                @input="updateDayOfMonth($event)"
              />
            </template>
            <template v-slot:hour>
              <ui-input
                v-model="hour"
                dense
                class="cronDigit"
                @input="updateHour($event)"
              />
            </template>
            <template v-slot:minute>
              <ui-input
                v-model="minute"
                dense
                class="cronDigit"
                @input="updateMinute($event)"
              />
            </template>
          </i18n>
        </span>
      </v-tab-item>
      <v-tab-item>
        <span class="tab">
          <i18n
            path="servers.WeeklyTab"
            class="tab"
          >
            <template v-slot:weekdays>
              <span
                style="display: inline-block;"
                class="px-2 py-0"
              >
                <v-checkbox
                  v-model="dayOfWeek"
                  dense
                  hide-details
                  :label="$t('common.Monday')"
                  value="1"
                />
                <v-checkbox
                  v-model="dayOfWeek"
                  dense
                  hide-details
                  :label="$t('common.Tuesday')"
                  value="2"
                />
                <v-checkbox
                  v-model="dayOfWeek"
                  dense
                  hide-details
                  :label="$t('common.Wednesday')"
                  value="3"
                />
                <v-checkbox
                  v-model="dayOfWeek"
                  dense
                  hide-details
                  :label="$t('common.Thursday')"
                  value="4"
                />
                <v-checkbox
                  v-model="dayOfWeek"
                  dense
                  hide-details
                  :label="$t('common.Friday')"
                  value="5"
                />
                <v-checkbox
                  v-model="dayOfWeek"
                  dense
                  hide-details
                  :label="$t('common.Saturday')"
                  value="6"
                />
                <v-checkbox
                  v-model="dayOfWeek"
                  dense
                  hide-details
                  :label="$t('common.Sunday')"
                  value="0"
                />
              </span>
            </template>
            <template v-slot:hour>
              <ui-input
                v-model="hour"
                dense
                class="cronDigit"
                @input="updateHour($event)"
              />
            </template>
            <template v-slot:minute>
              <ui-input
                v-model="minute"
                dense
                class="cronDigit"
                @input="updateMinute($event)"
              />
            </template>
          </i18n>
        </span>
      </v-tab-item>
      <v-tab-item>
        <span class="tab">
          <i18n path="servers.MonthlyTab">
            <template v-slot:monthInterval>
              <ui-input
                v-model="month"
                dense
                class="cronDigit"
                @input="updateMonth($event)"
              />
            </template>
            <template v-slot:day>
              <ui-input
                v-model="dayOfMonth"
                dense
                class="cronDigit"
                @input="updateDayOfMonth($event)"
              />
            </template>
            <template v-slot:hour>
              <ui-input
                v-model="hour"
                dense
                class="cronDigit"
                @input="updateHour($event)"
              />
            </template>
            <template v-slot:minute>
              <ui-input
                v-model="minute"
                dense
                class="cronDigit"
                @input="updateMinute($event)"
              />
            </template>
          </i18n>
        </span>
      </v-tab-item>
      <v-tab-item>
        <ui-input
          v-model="expression"
          class="px-2"
          :label="$t('servers.CronExpression')"
          @input="updateExpression($event)"
        />
      </v-tab-item>
    </v-tabs-items>
  </div>
</template>

<style scoped>
.tab {
  display: flex;
  align-items: center;
  justify-content: center;
}

.tab > * {
  padding: 0.25em;
}

.cronDigit {
  display: inline-block;
  max-width: 3em;
}
</style>

<script>
import cronParser from 'cron-parser'

const TAB_HOURLY = 0
const TAB_DAILY = 1
const TAB_WEEKLY = 2
const TAB_MONTHLY = 3
const TAB_ADVANCED = 4

export default {
  props: {
    value: { type: undefined, default: () => '', required: true }
  },
  data () {
    return {
      currentMode: 1,
      minute: '0',
      hour: '0',
      dayOfMonth: '*/1',
      month: '*',
      dayOfWeek: [0, 1, 2, 3, 4, 5, 6],
      expression: ''
    }
  },
  watch: {
    dayOfWeek (newVal) {
      this.updateDayOfWeek(newVal)
    },
    currentMode (newVal) {
      if (newVal === TAB_HOURLY && (this.hour === '0' || this.hour === '*')) {
        this.hour = '1'
      } else if (newVal === TAB_DAILY && (this.dayOfMonth === '0' || this.dayOfMonth === '*')) {
        this.dayOfMonth = '1'
      } else if (newVal === TAB_WEEKLY && this.dayOfWeek.length === 0) {
        this.dayOfWeek = [0]
      } else if (newVal === TAB_MONTHLY && (this.month === '0' || this.month === '*')) {
        this.month = '1'
      }

      this.$emit('input', this.buildExpression())
    }
  },
  mounted () {
    if (!this.value || this.value.trim() === '') {
      this.currentMode = 1
      this.minute = '0'
      this.hour = '0'
      this.dayOfMonth = '*/1'
      this.month = '*'
      this.dayOfWeek = [0, 1, 2, 3, 4, 5, 6]
    } else {
      this.expression = this.value
      const fields = cronParser.parseExpression(this.value).fields
      const ignoreDOM = this.value.split(' ')[2] === '?'
      const ignoreDOW = this.value.split(' ')[4] === '?'
      if (
        ignoreDOW &&
        fields.minute.length === 1 &&
        fields.hour.length === 1 &&
        fields.dayOfMonth.length === 1 &&
        fields.dayOfWeek.length === 8
      ) {
        this.currentMode = TAB_MONTHLY
      } else if (
        ignoreDOM &&
        fields.minute.length === 1 &&
        fields.hour.length === 1 &&
        fields.dayOfMonth.length === 31 &&
        fields.month.length === 12
      ) {
        this.currentMode = TAB_WEEKLY
      } else if (
        fields.minute.length === 1 &&
        fields.hour.length === 1 &&
        fields.month.length === 12 &&
        fields.dayOfWeek.length === 8
      ) {
        this.currentMode = TAB_DAILY
      } else if (
        fields.minute.length === 1 &&
        fields.dayOfMonth.length === 31 &&
        fields.month.length === 12 &&
        fields.dayOfWeek.length === 8
      ) {
        this.currentMode = TAB_HOURLY
      } else {
        this.currentMode = TAB_ADVANCED
      }

      if (this.mode !== 4) {
        const field = this.value.split(' ').map(e => {
          return e.replace(/^\*\//, '').replace(/^\*$/, '1')
        })
        this.minute = field[0]
        this.hour = field[1]
        this.dayOfMonth = field[2]
        this.month = field[3]
        this.dayOfWeek = fields.dayOfWeek.map(e => e.toString()).filter(e => e !== '7')
      }
    }
  },
  methods: {
    buildExpression () {
      if (this.currentMode === TAB_HOURLY) {
        this.expression = `${this.minute} */${this.hour} * * ?`
      } else if (this.currentMode === TAB_DAILY) {
        this.expression = `${this.minute} ${this.hour} */${this.dayOfMonth} * ?`
      } else if (this.currentMode === TAB_WEEKLY) {
        this.expression = `${this.minute} ${this.hour} ? * ${this.dayOfWeek.join(',')}`
      } else if (this.currentMode === TAB_MONTHLY) {
        this.expression = `${this.minute} ${this.hour} ${this.dayOfMonth} */${this.month} ?`
      }
      return this.expression
    },
    updateMinute (value) {
      this.minute = value
      this.$emit('input', this.buildExpression())
    },
    updateHour (value) {
      this.hour = value
      this.$emit('input', this.buildExpression())
    },
    updateDayOfMonth (value) {
      this.dayOfMonth = value
      this.$emit('input', this.buildExpression())
    },
    updateMonth (value) {
      this.month = value
      this.$emit('input', this.buildExpression())
    },
    updateDayOfWeek (value) {
      this.dayOfWeek = value
      this.$emit('input', this.buildExpression())
    },
    updateExpression (expr) {
      this.expression = expr
      this.$emit('input', expr)
    }
  }
}
</script>
