<template>
  <v-layout>
    <v-col cols="3" id="mailList">
      <v-list>
        <v-list-item v-for="mail in mails"></v-list-item>
      </v-list>
    </v-col>
    <v-col cols="9" id="content">

    </v-col>
  </v-layout>
</template>

<script>
export default {
  name: "Home",
  data: () => {
    return {
      receivedMails: []
    }
  },
  computed: {
    mails() {
      return this.receivedMails
    }
  }
}
</script>

<style scoped>
#mailList {
  height: 100%;
  background: blue;
}

#content {
  height: 100%;
  background: red;
}
</style>