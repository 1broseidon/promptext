import React from 'react';
import ComponentCreator from '@docusaurus/ComponentCreator';

export default [
  {
    path: '/promptext/__docusaurus/debug',
    component: ComponentCreator('/promptext/__docusaurus/debug', 'e5b'),
    exact: true
  },
  {
    path: '/promptext/__docusaurus/debug/config',
    component: ComponentCreator('/promptext/__docusaurus/debug/config', 'ae8'),
    exact: true
  },
  {
    path: '/promptext/__docusaurus/debug/content',
    component: ComponentCreator('/promptext/__docusaurus/debug/content', '28f'),
    exact: true
  },
  {
    path: '/promptext/__docusaurus/debug/globalData',
    component: ComponentCreator('/promptext/__docusaurus/debug/globalData', 'd39'),
    exact: true
  },
  {
    path: '/promptext/__docusaurus/debug/metadata',
    component: ComponentCreator('/promptext/__docusaurus/debug/metadata', '681'),
    exact: true
  },
  {
    path: '/promptext/__docusaurus/debug/registry',
    component: ComponentCreator('/promptext/__docusaurus/debug/registry', '653'),
    exact: true
  },
  {
    path: '/promptext/__docusaurus/debug/routes',
    component: ComponentCreator('/promptext/__docusaurus/debug/routes', 'ce2'),
    exact: true
  },
  {
    path: '/promptext/',
    component: ComponentCreator('/promptext/', 'c0b'),
    routes: [
      {
        path: '/promptext/',
        component: ComponentCreator('/promptext/', 'b83'),
        routes: [
          {
            path: '/promptext/',
            component: ComponentCreator('/promptext/', 'a10'),
            routes: [
              {
                path: '/promptext/configuration',
                component: ComponentCreator('/promptext/configuration', 'cf8'),
                exact: true,
                sidebar: "tutorialSidebar"
              },
              {
                path: '/promptext/file-filtering',
                component: ComponentCreator('/promptext/file-filtering', 'f16'),
                exact: true,
                sidebar: "tutorialSidebar"
              },
              {
                path: '/promptext/getting-started',
                component: ComponentCreator('/promptext/getting-started', '4c5'),
                exact: true,
                sidebar: "tutorialSidebar"
              },
              {
                path: '/promptext/output-formats',
                component: ComponentCreator('/promptext/output-formats', '1bd'),
                exact: true,
                sidebar: "tutorialSidebar"
              },
              {
                path: '/promptext/performance',
                component: ComponentCreator('/promptext/performance', '578'),
                exact: true,
                sidebar: "tutorialSidebar"
              },
              {
                path: '/promptext/project-analysis',
                component: ComponentCreator('/promptext/project-analysis', '0bb'),
                exact: true,
                sidebar: "tutorialSidebar"
              },
              {
                path: '/promptext/token-analysis',
                component: ComponentCreator('/promptext/token-analysis', '3f0'),
                exact: true,
                sidebar: "tutorialSidebar"
              },
              {
                path: '/promptext/',
                component: ComponentCreator('/promptext/', 'f2a'),
                exact: true,
                sidebar: "tutorialSidebar"
              }
            ]
          }
        ]
      }
    ]
  },
  {
    path: '*',
    component: ComponentCreator('*'),
  },
];
