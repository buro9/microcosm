// This file handles displaying an event page, rather than the create/edit of an event

/////////
// MAP //
/////////
(function(){

  var EventMap = (function(){
    var eventMap = function(options){
      // reference to map marker
      this.marker = [];

      var mapID = 'map-canvas';
      if (typeof options.map !== 'undefined'){
        mapID = options.map;
      }

      this.map = L.map( options.map , { zoomAnimation : false })
                  .setMaxBounds([[-90,-180],[90,180]]); // Restrict map to valid lat:lng pairs

      var googleLayer = new L.Google('ROADMAP');
      this.map.addLayer(googleLayer);

      var attribution = {
        attribution: 'Maps &copy; <a href="http://www.thunderforest.com">Thunderforest</a>, Data &copy; <a href="http://www.openstreetmap.org/copyright">OpenStreetMap contributors</a>',
        maxZoom: 22
      }

      var osmCycleLayer = new L.TileLayer("https://a.tile.thunderforest.com/cycle/{z}/{x}/{y}.png", attribution);
      var osmLandscapeLayer = new L.TileLayer("https://a.tile.thunderforest.com/landscape/{z}/{x}/{y}.png", attribution);

      this.map.addControl(new L.Control.Layers({
        'Open Cycle Map' : osmCycleLayer,
        'Landscape'      : osmLandscapeLayer,
        'Google Maps'    : googleLayer
        })
      );

      return this;
    };

    eventMap.prototype.addMarker = function(latlng){
      var newMarker = L.marker(latlng);
      this.marker.push(newMarker);
      return this;
    };

    eventMap.prototype.renderMarkers = function(){
      if (this.marker.length > 0){
        for(var i=0,j=this.marker.length;i<j;i++){
          this.map.addLayer( this.marker[i] );
        }
      }
      return this;
    };

    // If this is the edit screen and we need to set the map to known values then
    // we do so here
    // @param lat number
    // @param lng number
    // @param bounds array(array(north, west),array(south, east))
    eventMap.prototype.restore = function(lat, lng, bounds){

      this.map.fitBounds(bounds);
      // Bounds will always be 1 level too far out due to imprecision in the
      // numbers and Leaflet aggressively ensuring the bounds fit inside the
      // map area. To prevent us slowly zooming out with each edit, we zoom in
      this.map.zoomIn();

      this.addMarker(new L.LatLng(lat, lng));
      this.renderMarkers();
    };

    return eventMap;

  })();

  window.EventMap = EventMap;

})();