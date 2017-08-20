Vue.component('album', {
    template: '#album-col',
    props : {
        album: Object
    },
    data: function () {
        return {
            checkedTracks: [],
            show: false,
            checked: true
        }
    },
    computed:  {
        allTracksChecked: function () {
            return this.checkedTracks.length === this.album.tracks.items.length
        }
    },
    created: function() {
        var self = this;
        self.album.tracks.items.forEach(function (track) {
            self.checkedTracks.push(track.id)
        })
    },
    methods: {
        updateAlbum: function () {
            var self = this;
            if(self.allTracksChecked) {
                self.checked = true;
            }
            self.album.tracks.items.forEach(function (track) {
                var index = self.checkedTracks.indexOf(track.id);
                if(self.checked) {
                    if(index < 0) {
                        self.checkedTracks.push(track.id);
                        self.updateTrack(track.id)
                    }
                }
                else {
                    if(index > -1) {
                        self.checkedTracks.splice(index, 1);
                        self.updateTrack(track.id)
                    }
                }
            });
        },
        updateTrack: function(trackId) {
            this.$emit('update', trackId);
        }
    }
});