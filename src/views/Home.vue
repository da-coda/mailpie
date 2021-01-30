<template>
  <v-layout>
   <v-img src="@/assets/img/logo.png" max-width="400px" max-height="400px"/>
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