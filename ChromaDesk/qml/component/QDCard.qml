// Fluent Design Card Component
import QtQuick

Rectangle {
    id: control
    
    // ============ Custom Properties ============
    
    property bool hoverable: false
    property bool clickable: false
    property int elevation: 1  // 1-3, higher = more shadow
    
    signal clicked()
    
    // ============ Style ============
    
    radius: Theme.radiusLarge
    color: mouseArea.containsMouse && hoverable ? Theme.surfaceHover : Theme.surface
    border.width: Theme.borderWidthThin
    border.color: Theme.border
    
    scale: mouseArea.containsMouse && hoverable ? 1.02 : 1.0
    
    Behavior on color {
        ColorAnimation {
            duration: Theme.animationDurationFast
            easing.type: Theme.animationEasingType
        }
    }
    
    Behavior on scale {
        NumberAnimation {
            duration: Theme.animationDurationMedium
            easing.type: Easing.OutCubic
        }
    }
    
    // ============ Mouse Area ============
    
    MouseArea {
        id: mouseArea
        anchors.fill: parent
        hoverEnabled: control.hoverable || control.clickable
        cursorShape: control.clickable ? Qt.PointingHandCursor : Qt.ArrowCursor
        
        onClicked: {
            if (control.clickable) {
                control.clicked()
            }
        }
    }
}
