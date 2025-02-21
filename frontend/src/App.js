import React from 'react';
import { BrowserRouter as Router, Route, Routes } from 'react-router-dom';
import Home from './components/courses/Home';
import Login from './components/auth/Login';
import Register from './components/auth/Register';
import CourseDetails from './components/courses/CourseDetails';
import MyCourses from './components/courses/MyCourses';
import UploadFile from './components/courses/UploadFile';
import CourseFiles from './components/courses/CourseFiles';
import SearchCourses from './components/courses/SearchCourses';
import ProtectedRoute from './components/ProtectedRoute';
import AdminRoute from './components/AdminRoute';
import ManageCourses from './components/courses/ManageCourses';
import AddCourse from './components/courses/AddCourse';
import EditCourse from './components/courses/EditCourse';
import CourseComments from './components/courses/CourseComments'; 
import CommentForm from './components/courses/CommentForm';
import { Navigate } from 'react-router-dom';

function App() {
  return (
    <Router>
      <Routes>
        <Route path="/" element={<Navigate replace to="/login" />} />
        <Route path="/login" element={<Login />} />
        <Route path="/register" element={<Register />} />
        
        <Route element={<ProtectedRoute />}>
          <Route path="/home" element={<Home />} />
          <Route path="/courses/:courseId" element={<CourseDetails />} />
          <Route path="/my-courses" element={<MyCourses />} />
          <Route path="/search" element={<SearchCourses />} />
          <Route path="/upload/:courseId" element={<UploadFile />} />
          <Route path="/courses/:courseId/comments" element={<CourseComments />} />
          <Route path="/courses/:courseId/files" element={<CourseFiles />} />
          <Route path="/courses/:courseId/comment" element={<CommentForm />} />
        </Route>

        <Route element={<AdminRoute />}>
          <Route path="/manage-courses" element={<ManageCourses />} />
          <Route path="/add-course" element={<AddCourse />} />
          <Route path="/edit-course/:courseId" element={<EditCourse />} />
        </Route>
      </Routes>
    </Router>
  );
}

export default App;
