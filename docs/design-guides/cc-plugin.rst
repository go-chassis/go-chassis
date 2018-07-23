Config Client Plugin
====================

Config Client
#############
Go-Chassis provides the functionality to pull the configs from different config-centers, to keep the go-chassis extensible to support multiple config-center this Client was implemented as a plugin.

More Details of the plugin can be found  `here <../dev-guides/archaius-config-source-plugin>`_

Currently Go-Chassis has two implementation of this plugin for Huawei Config-Center and Ctrip Apollo Config-Center.

Basic Sequence diagram for this plugin is given below.

 .. image:: images/CC-Plugin.png
    :alt: CC Plugin Sequence Diagram