// Zoom persistence utility - DISABLED to fix zoom issues
// This script was causing problems with browser zoom behavior
(function() {
  'use strict';
  
  // Clear any previously stored zoom levels that might be causing issues
  try {
    localStorage.removeItem('browserZoomLevel');
  } catch (error) {
    console.log('Could not clear zoom level:', error);
  }
})();
