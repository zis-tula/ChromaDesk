// Fluent Design Divider Component
import QtQuick

Rectangle {
    id: control
    
    // ============ Custom Properties ============
    
    property bool vertical: false
    property string label: ""
    
    // ============ Size & Style ============
    
    implicitWidth: vertical ? 1 : 200
    implicitHeight: vertical ? 200 : (label !== "" ? labelText.height + Theme.spacingSmall * 2 : 1)
    
    color: label !== "" ? "transparent" : Theme.border
    
    // ============ Horizontal divider with label ============
    
    Row {
        visible: !vertical && label !== ""
        anchors.fill: parent
        spacing: Theme.spacingMedium
        
        Rectangle {
            width: (parent.width - labelText.width - Theme.spacingMedium * 2) / 2
            height: 1
            color: Theme.border
            anchors.verticalCenter: parent.verticalCenter
        }
        
        Text {
            id: labelText
            text: control.label
            font.family: Theme.fontFamily
            font.pixelSize: Theme.fontSizeSmall
            color: Theme.textSecondary
            anchors.verticalCenter: parent.verticalCenter
        }
        
        Rectangle {
            width: (parent.width - labelText.width - Theme.spacingMedium * 2) / 2
            height: 1
            color: Theme.border
            anchors.verticalCenter: parent.verticalCenter
        }
    }
}
