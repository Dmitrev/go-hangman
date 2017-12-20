new Vue({
  el: '#app',

  data: {
      ws: null, // websocket
      word:[],
      lives: 0,
      letter: null
  },

  methods: {
    send: function(){
      if( !this.letter.match(/[A-Za-z]{1}/) ){
        alert("INVALID input");
        this.letter = "";
      }

      this.ws.send(JSON.stringify({
        type: "turn",
        data: this.letter
      }));

      this.letter = "";
    }
  },

  created: function() {
    var self = this;
    this.ws = new WebSocket('ws://' + window.location.host + '/ws');

    this.ws.onmessage = function(event){
      var json = JSON.parse(event.data)
      self.word = json.word;
      self.lives = json.lives;
    };

  }

});
