/*
 * L.TileLayer is used for standard xyz-numbered tile layers.
 */
L.Google = L.Class.extend({
	includes: L.Mixin.Events,

	options: {
		minZoom: 0,
		maxZoom: 18,
		tileSize: 256,
		subdomains: 'abc',
		errorTileUrl: '',
		attribution: '',
		opacity: 1,
		continuousWorld: false,
		noWrap: false,
	},

	// Possible types: SATELLITE, ROADMAP, HYBRID
	initialize: function(type, options) {
		L.Util.setOptions(this, options);

		this._type = google.maps.MapTypeId[type || 'SATELLITE'];
	},

	onAdd: function(map, insertAtTheBottom) {
		this._map = map;
		this._insertAtTheBottom = insertAtTheBottom;

		// create a container div for tiles
		this._initContainer();
		this._initMapObject();

		// set up events
		map.on('viewreset', this._resetCallback, this);

		this._limitedUpdate = L.Util.limitExecByInterval(this._update, 150, this);
		map.on('move', this._update, this);
		//map.on('moveend', this._update, this);

		this._reset();
		this._update();
	},

	onRemove: function(map) {
		this._map._container.removeChild(this._container);
		//this._container = null;

		this._map.off('viewreset', this._resetCallback, this);

		this._map.off('move', this._update, this);
		//this._map.off('moveend', this._update, this);
	},

	getAttribution: function() {
		return this.options.attribution;
	},

	setOpacity: function(opacity) {
		this.options.opacity = opacity;
		if (opacity < 1) {
			L.DomUtil.setOpacity(this._container, opacity);
		}
	},

	_initContainer: function() {
		var tilePane = this._map._container
			first = tilePane.firstChild;

		if (!this._container) {
			this._container = L.DomUtil.create('div', 'leaflet-google-layer leaflet-top leaflet-left');
			this._container.id = "_GMapContainer";
		}

		if (true) {
			tilePane.insertBefore(this._container, first);

			this.setOpacity(this.options.opacity);
			var size = this._map.getSize();
			this._container.style.width = size.x + 'px';
			this._container.style.height = size.y + 'px';
			this._container.style.zIndex = 0;
		}
	},

	_initMapObject: function() {
		this._google_center = new google.maps.LatLng(0, 0);
		var map = new google.maps.Map(this._container, {
		    center: this._google_center,
		    zoom: 0,
		    mapTypeId: this._type,
		    disableDefaultUI: true,
		    keyboardShortcuts: false,
		    draggable: false,
		    disableDoubleClickZoom: true,
		    scrollwheel: false,
		    streetViewControl: false
		});

		var _this = this;
		this._reposition = google.maps.event.addListenerOnce(map, "center_changed", 
			function() { _this.onReposition(); });
	
		map.backgroundColor = '#ff0000';
		this._google = map;
	},

	_resetCallback: function(e) {
		this._reset(e.hard);
	},

	_reset: function(clearOldContainer) {
		this._initContainer();
	},

	_update: function() {
		this._resize();

		var bounds = this._map.getBounds();
		var ne = bounds.getNorthEast();
		var sw = bounds.getSouthWest();
		var google_bounds = new google.maps.LatLngBounds(
			new google.maps.LatLng(sw.lat, sw.lng),
			new google.maps.LatLng(ne.lat, ne.lng)
		);
		var center = this._map.getCenter();
		var _center = new google.maps.LatLng(center.lat, center.lng);

		this._google.setCenter(_center);
		this._google.setZoom(this._map.getZoom());
		//this._google.fitBounds(google_bounds);
	},

	_resize: function() {
		var size = this._map.getSize();
		if (this._container.style.width == size.x &&
		    this._container.style.height == size.y)
			return;
		this._container.style.width = size.x + 'px';
		this._container.style.height = size.y + 'px';
		google.maps.event.trigger(this._google, "resize");
	},

	onReposition: function() {
		//google.maps.event.trigger(this._google, "resize");
	}
});

// Create map
var map = L.map('map', {zoomAnimation:false})
	.setMaxBounds([[-90,-180],[90,180]]); // Restrict map to valid lat:lng pairs

var googleLayer = new L.Google('ROADMAP');
map.addLayer(googleLayer);

//var cloudmadeLayer = new L.TileLayer("https://d1qte70nkdppk5.cloudfront.net/d6f1a0c60e9746faa7cbfaec4b92dff3/96931/256/{z}/{x}/{y}.png");
var osmLayer = new L.TileLayer("http://otile1.mqcdn.com/tiles/1.0.0/map/{z}/{x}/{y}.jpg");
map.addControl(new L.Control.Layers({'Open Street Map':osmLayer, 'Google Maps': googleLayer}))

// If this is the edit screen and we need to set the map to known values then
// we do so here
function restoreLocationState(lat, lng, bounds) {
	map.fitBounds(bounds);

	// Bounds will always be 1 level too far out due to imprecision in the
	// numbers and Leaflet aggressively ensuring the bounds fit inside the
	// map area. To prevent us slowly zooming out, we zoom in
	map.zoomIn();

	var marker = L.marker(new L.LatLng(lat, lng));
	map.addLayer(marker);

	if (lat != 0 && lng != 0 && bounds.length == 2 && bounds[1].length == 2) {

		var ua = navigator.userAgent.toLowerCase();
		var isAndroid = ua.indexOf("android") > -1;

		var gmaplink = '';
		if (isAndroid) {
			gmaplink = 'geo:' + lat + ',' + lng + '?q=' + lat + ',' + lng + '(' + $('#where').text() + ')';
		} else {
			gmaplink = 'https://maps.google.com/maps?q=' + lat + ',' + lng + 
			'(' + $('#where').text() + ')&spn=' + (bounds[0][0]-bounds[1][0]) + ',' + (bounds[0][1]-bounds[1][1]);
		}

		var osmlink = 'http://www.openstreetmap.org/?minlon=' + bounds[0][1] + '&minlat=' + bounds[1][0] +
			'&maxlon=' + bounds[1][1] + '&maxlat=' + bounds[0][0] +
			'&box=no&mlat=' + lat + '&mlon=' + lng;

		if (isAndroid) {
			$('#maplinks').append(
				'View location in <a href="' + gmaplink + '">Maps</a>.'
			).show();
		} else {
			$('#maplinks').append(
				'View location in <a href="' + gmaplink + '">Google Maps</a> or <a href="' + osmlink + '">Open Street Maps</a>.'
			).show();
		}
	}
}

function formatDate(d) {
	var monthsShort = ["Jan", "Feb", "Mar", "Apr", "May", "Jun", "Jul", "Aug", "Sep", "Oct", "Nov", "Dec"];

	var hh = d.getHours();
	var pm = (hh > 12);

	return d.getDate() + ' ' +
		monthsShort[d.getMonth()] + ' ' + 
		d.getFullYear() + ', ' +
		(hh > 12 ? hh - 12: hh) + ':' +
		d.getMinutes() + ' ' +
		(pm ? 'p.m.' : 'a.m.');
}

function getDateFromIsoString(input) {
	var parts = input.match(/(\d+)/g);
	// new Date(year, month [, date [, hours[, minutes[, seconds[, ms]]]]])
	return new Date(parts[0], parts[1]-1, parts[2], parts[3], parts[4]); // months are 0-based
}

function showEndDate(startDate, duration) {
	// We need a date object
	var sd = getDateFromIsoString(startDate);
	if (sd) {
		// And the duration in minutes
		dur = Number(duration);
		if (dur) {
			// To figure out the millisecond difference
			dur = dur * 60000;
			// And create a date object representing that end point in time
			var ed = moment(new Date(sd.valueOf() + dur));
			// And update the fields for the end info
			$('#ends').attr('datetime', ed.utc().format());
			// This is in base.js
			updateTimes();
		}
	}
}
