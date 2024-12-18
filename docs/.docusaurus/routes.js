import React from 'react';
import ComponentCreator from '@docusaurus/ComponentCreator';

export default [
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
