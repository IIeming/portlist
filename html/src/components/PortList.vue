<template>
  <div class="portlist-windwos">
    <div class="portlist-show">
      <div>
        <h1>可用端口<span>&nbsp;:</span></h1>
        <span class="font-style">{{ data }}</span>
      </div>
      <div class="portlist-input" v-if="portListInput">
        <input
          type="text"
          maxlength="1"
          :key="0"
          @input="handleInput(portOne, 0)"
          v-model="portOne"
          @blur="validateNumber(portOne)"
        />
        <input
          type="text"
          maxlength="1"
          :key="1"
          @input="handleInput(portTwo, 1)"
          v-model="portTwo"
          @blur="validateNumber(portTwo)"
        />
        <input
          type="text"
          maxlength="1"
          :key="2"
          @input="handleInput(portThree, 2)"
          v-model="portThree"
          @blur="validateNumber(portThree)"
        />
        <input
          type="text"
          :key="3"
          @input="handleInput(portFour, 3)"
          maxlength="1"
          v-model="portFour"
          @blur="validateNumber(portFour)"
        />
        <input
          type="text"
          :key="4"
          maxlength="1"
          v-model="portFive"
          @blur="validateNumber(portFive)"
        />
        <button @click="portReset">确认</button>
      </div>

      <div class="portlist-label">
        <label v-if="portOneClue">请输入1-6范围内数字</label>
        <label v-if="portNumberClue">请输入0-9范围内数字</label>
        <label v-if="portButtonClue">请输入1-65535范围内数字</label>
        <label v-if="portSucceeClue">端口修改失败，请联系管理员</label>
      </div>
      <div class="portlist-button" v-if="!portListInput">
        <button @click="portRestShow">修改</button>
        <button @click="portUsed">已使用</button>
      </div>
    </div>
    <div class="portlist-plus">
      <div class="plus-current">
        <h1>今天是</h1>
        <span>{{ dateWeekDay }}</span>
      </div>
      <div class="plus-holiday">
        <h1>距{{ dataFestival }}假期</h1>
        <p class="plus-font-style">
          还有<span>{{ dataHoliday }}天</span>
        </p>
      </div>
    </div>
  </div>
</template>

<script>
import { ref, onMounted, nextTick } from "vue";
import request from "../request.js";

export default {
  name: "",
  setup() {
    const data = ref("NULL");
    // 输入框输入内容
    const portOne = ref("");
    const portTwo = ref("");
    const portThree = ref("");
    const portFour = ref("");
    const portFive = ref("");
    const index = ref(0);

    // 标签显示
    const portOneClue = ref(false);
    const portNumberClue = ref(false);
    const portButtonClue = ref(false);
    const portSucceeClue = ref(false);
    const portListInput = ref(false);

    // 日期
    const dateWeekDay = ref("星期一");
    const dataHoliday = ref("365");
    const dataFestival = ref("下一个");

    onMounted(async () => {
      portSucceeClue.value = false;
      try {
        await request({
          url: "get",
          headers: {
            "Content-Type": "application/json",
          },
        }).then((response) => {
          // alert('请求成功');
          // console.log(JSON.stringify(response.data));
          data.value = response.data.value;
        });
      } catch (error) {
        console.error(error);
      }
    });
    // 获取日期
    request({
      url: "date",
      headers: {
        "Content-Type": "application/json",
      },
    })
      .then((response) => {
        // alert('请求成功');
        // console.log(JSON.stringify(response.data));
        dataHoliday.value = response.data.value;
        dateWeekDay.value = response.data.weekday;
        dataFestival.value = response.data.festival;
      })
      .catch((error) => {
        console.log(error.response.status, error.response.data);
      });

    function portUsed() {
      request({
        url: "add",
        headers: {
          "Content-Type": "application/json",
        },
      }).then((response) => {
        // alert('请求成功');
        // console.log(JSON.stringify(response.data));
        data.value = response.data.value;
      });
    }

    function portReset() {
      const value =
        portOne.value +
        portTwo.value +
        portThree.value +
        portFour.value +
        portFive.value;
      const port = parseInt(value);
      if (value === "") {
        portListInput.value = false;
      } else {
        // console.log("port:", port, value, isNaN(value));
        if (!isNaN(value) && port > 0 && port < 65536) {
          // console.log("yesport:", port);
          portButtonClue.value = false;
          // 请求reset接口
          const formData = new FormData();
          formData.append("curvePost", port);
          // console.log(formData);
          request({
            url: "reset",
            data: formData,
            headers: {
              "Content-Type": "multipart/form-data",
            },
          })
            .then((response) => {
              // alert('请求成功');
              data.value = response.data.value;
              portSucceeClue.value = false;
              portListInput.value = false;
            })
            .catch((error) => {
              console.log(error.response.status, error.response.data);
              portSucceeClue.value = true;
            });
        } else {
          // console.log("noport:", port);
          portButtonClue.value = true;
        }
      }
      portOne.value = "";
      portTwo.value = "";
      portThree.value = "";
      portFour.value = "";
      portFive.value = "";
    }

    function portRestShow() {
      portListInput.value = true;
    }

    // 校验第一个输入框输入范围是否在 1-6
    function validateInput() {
      const value = parseInt(portOne.value);
      // console.log(portOne.value, value);
      if (isNaN(value) || value > 6 || value < 1) {
        // 输入不在0-9的范围内，执行相应的操作
        portOneClue.value = true;
      } else {
        portOneClue.value = false;
      }
    }

    // 检验输入的是否是数字
    function validateNumber(number) {
      // const value = parseInt(number);
      // console.log("number", number);
      if (number.trim() !== "" && isNaN(number)) {
        portNumberClue.value = true;
      } else {
        portNumberClue.value = false;
      }
    }

    // 自动聚焦下一个输入框
    function handleInput(number, value) {
      index.value = value + 1;
      if (number.length === 1) {
        // console.log("index", index.value);
        nextTick(() => {
          const nextInput = document.querySelectorAll("input")[index.value];
          nextInput.focus();
        });
      }
    }

    return {
      data,
      portOne,
      portTwo,
      portThree,
      portFour,
      portFive,
      portOneClue,
      portNumberClue,
      portButtonClue,
      portSucceeClue,
      portListInput,
      dataHoliday,
      dateWeekDay,
      dataFestival,
      portUsed,
      portReset,
      portRestShow,
      validateInput,
      validateNumber,
      handleInput,
    };
  },
};
</script>

<style scoped>
.portlist-windwos {
  height: 400px;
  width: 600px;
  background: linear-gradient(to right, rgb(0, 225, 255), rgb(0, 132, 255));
  box-shadow: 0 0 10px rgba(0, 0, 0, 0.5);
  background-size: 600px 400px;
  position: absolute;
  align-items: center;
  display: flex;
  flex-direction: row;
  justify-content: center;
  left: 50%;
  top: 50%;
  transform: translate(-50%, -50%);
  border-radius: 10px;
}

.font-style {
  font-size: 3rem;
  color: white;
  font-style: normal;
  font-weight: 600;
}

.portlist-show {
  height: 100%;
  width: 55%;
  /* border: 2px solid black; */
}

.portlist-show div {
  margin-left: 3rem;
  margin-right: 3rem;
}

.portlist-show h1 {
  color: black;
  margin-bottom: 0.3rem;
  margin-top: 3rem;
  font-size: 1.6rem;
}

.portlist-show h1 span {
  color: white;
}

.portlist-input {
  margin-top: 9.5rem;
  display: flex;
  flex-direction: row;
  justify-content: space-between;
  align-items: center;
}

.portlist-input input {
  outline: none;
  width: 1.5rem;
  height: 1.5rem;
  font-size: 1.2rem;
  font-weight: 500;
  text-align: center;
  border: 1px solid white;
}

.portlist-input input:focus {
  padding: 0px 1px;
  border: 2px solid #00000080;
  border-radius: 1px;
}

.portlist-input button {
  margin-left: 0.8rem;
  background-color: white;
  font-size: 0.9rem;
  color: rgb(0 212 255 / 91%);
  border: 0px solid white;
  border-radius: 3px;
  font-weight: 600;
  padding: 0px 1rem;
  height: 1.8rem;
}

.portlist-input button:active {
  background-color: #dadce0;
  border: 0px solid #dadce0;
}

.portlist-label label {
  outline: none;
  color: red;
  font-size: 0.5rem;
  font-weight: 600;
  letter-spacing: 0px;
}

.portlist-button {
  margin-top: 9.5rem;
  display: flex;
  flex-direction: row;
  justify-content: space-between;
}

.portlist-button button {
  width: 6rem;
  height: 2.5rem;
  background-color: white;
  font-size: 0.9rem;
  color: rgb(0 212 255 / 91%);
  border: 1px solid white;
  border-radius: 3px;
  font-weight: 600;
}

.portlist-button button:active {
  /* color: white; */
  background-color: #dadce0;
  border: 1px solid #dadce0;
}

.portlist-plus {
  height: 100%;
  width: 45%;
  /* color: aqua; */
}

.portlist-plus div {
  margin-left: 2rem;
  margin-right: 3rem;
}

.portlist-plus h1 {
  color: rgb(27, 25, 25);
  margin-bottom: 0.3rem;
  margin-top: 3rem;
  font-size: 1.6rem;
}

.portlist-plus span {
  color: black;
  font-size: 3rem;
  color: white;
  font-style: normal;
  font-weight: 600;
}

.portlist-plus p {
  outline: none;
  margin-top: 0rem;
  color: white;
  margin-bottom: 0.3rem;
  font-size: 1.6rem;
  font-weight: bold;
}

.plus-current {
  margin-top: 3rem;
}

.plus-holiday {
  margin-top: 3rem;
  display: flex;
  flex-direction: column;
  align-items: flex-start;
}

/* 媒体查询 */
@media screen and (max-width: 620px) {
  .portlist-windwos {
    height: 400px;
    width: 80%;
    background: linear-gradient(to right, rgb(0, 225, 255), rgb(0, 132, 255));
    box-shadow: 0 0 10px rgba(0, 0, 0, 0.5);
    background-size: 600px 400px;
    position: absolute;
    align-items: center;
    display: flex;
    flex-direction: row;
    justify-content: center;
    left: 50%;
    top: 50%;
    transform: translate(-50%, -50%);
    border-radius: 10px;
  }
}
</style>